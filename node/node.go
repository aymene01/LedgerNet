package node

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/aymene01/ledgerNet/crypto"
	"github.com/aymene01/ledgerNet/pb"
	"github.com/aymene01/ledgerNet/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

const blockTime = time.Second * 5

type Mempool struct {
	txx map[string]*pb.Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: make(map[string]*pb.Transaction),
	}
}

func (pool *Mempool) Has(tx *pb.Transaction) bool {
	hash := hex.EncodeToString(types.HashTransaction(tx))
	_, ok := pool.txx[hash]
	return ok
}

func (pool *Mempool) Add(tx *pb.Transaction) bool {
	if pool.Has(tx) {
		return false
	}

	hash := hex.EncodeToString(types.HashTransaction(tx))
	pool.txx[hash] = tx
	return true
}

type ServerConfig struct {
	Version    string
	ListenAddr string
	PrivateKey *crypto.PrivateKey
}

type Node struct {
	ServerConfig
	logger zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[pb.NodeClient]*pb.Version

	mempool *Mempool
	pb.UnimplementedNodeServer
}

func NewNode(cfg ServerConfig) *Node {
	logger, _ := getLoggerConfig()
	return &Node{
		peers:        make(map[pb.NodeClient]*pb.Version),
		logger:       *logger.Sugar(),
		mempool:      NewMempool(),
		ServerConfig: cfg,
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.ListenAddr = listenAddr

	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)

	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return err
	}

	pb.RegisterNodeServer(grpcServer, n)
	n.logger.Infow("node started ...", "port", n.ListenAddr)

	// bootstrapNode
	if len(bootstrapNodes) > 0 {
		go n.bootstrapNetwork(bootstrapNodes)
	}

	if n.PrivateKey != nil {
		go n.validatorLoop()
	}

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *pb.Version) (*pb.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *pb.Transaction) (*pb.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	hash := hex.EncodeToString(types.HashTransaction(tx))

	if n.mempool.Add(tx) {
		n.logger.Debugw("received tx:", "from", peer.Addr, "hash", hash, "we", n.ListenAddr)
		go func() {
			if err := n.broadcast(tx); err != nil {
				n.logger.Errorw("broadcast error", "err", err)
			}
		}()
	}
	return &pb.Ack{}, nil
}

func (n *Node) getVersion() *pb.Version {
	return &pb.Version{
		Version:    "blocker-0.1",
		Height:     0,
		ListenAddr: n.ListenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) addPeer(c pb.NodeClient, v *pb.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// Handle the logic where we decide to
	// accept or drop the incomming node

	n.peers[c] = v

	// connect to the list of peers in the received list of peers

	if len(v.PeerList) > 0 {
		go n.bootstrapNetwork(v.PeerList)
	}
	n.logger.Debugw("new peer successfully connected",
		"we", n.ListenAddr,
		"remoteNode", v.ListenAddr,
		"height", v.Height,
	)
}

func (n *Node) deletePeer(c pb.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, c)
}

func (n *Node) broadcast(msg any) error {
	for peer := range n.peers {
		switch v := msg.(type) {
		case *pb.Transaction:
			_, err := peer.HandleTransaction(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *Node) bootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}

		n.logger.Debugw("dialing remote node", "we", n.ListenAddr, "remote", addr)
		c, v, err := n.newRemoteClient(addr)
		if err != nil {
			return err
		}

		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) validatorLoop() {
	n.logger.Infow("starting validator loop", "pubkey", n.PrivateKey.Public(), "blockTime", blockTime)
	ticker := time.NewTicker(blockTime)
	for {
		<-ticker.C
		// usdafe code ⚠️
		// for hash := range n.mempool.txx {
		// 	delete(n.mempool.txx, hash)
		// }
		n.logger.Debugw("time to create a new block", "lenTx", len(n.mempool.txx))
	}
}

func (n *Node) newRemoteClient(addr string) (pb.NodeClient, *pb.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Handshake(context.Background(), n.getVersion())
	if err != nil {
		return nil, nil, err
	}

	return c, v, nil
}

func (n *Node) canConnectWith(addr string) bool {
	if n.ListenAddr == addr {
		return false
	}

	connectedPeers := n.getPeerList()

	for _, connectAddr := range connectedPeers {
		if addr == connectAddr {
			return false
		}
	}

	return true
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := []string{}
	for _, version := range n.peers {
		peers = append(peers, version.ListenAddr)
	}

	return peers
}

func makeNodeClient(listenAddr string) (pb.NodeClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	c, err := grpc.NewClient(listenAddr, opts...)

	if err != nil {
		return nil, err
	}

	return pb.NewNodeClient(c), nil
}

func getPortNum(listenAddr string) (string, error) {
	values := strings.Split(listenAddr, ":")
	if len(values) != 2 {
		return "", errors.New("invalid listen value")
	}

	return values[1], nil
}

func (n *Node) dialRemoteNode(addr string) (pb.NodeClient, *pb.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Handshake(context.Background(), n.getVersion())
	if err != nil {
		return nil, nil, err
	}

	return c, v, nil
}

func getLoggerConfig() (*zap.Logger, error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	loggerConfig.EncoderConfig.ConsoleSeparator = " "
	loggerConfig.EncoderConfig.TimeKey = ""
	loggerConfig.DisableStacktrace = true
	return loggerConfig.Build()
}
