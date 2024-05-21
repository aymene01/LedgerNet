package node

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/aymene01/ledgerNet/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string
	logger     zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[pb.NodeClient]*pb.Version

	pb.UnimplementedNodeServer
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()

	return &Node{
		peers:   make(map[pb.NodeClient]*pb.Version),
		version: "blocker-1",
		logger:  *logger.Sugar(),
	}
}

func (n *Node) addPeer(c pb.NodeClient, v *pb.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Debugf("[%s] new peer connected [%s]", n.listenAddr, v.ListenAddr)
	n.peers[c] = v
}

func (n *Node) deletePeer(c pb.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, c)
}

func (n *Node) getVersion() *pb.Version {
	return &pb.Version{
		Version:    "blocker-0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Errorf("Handshake err", err)
			continue
		}

		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr

	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)

	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return err
	}

	pb.RegisterNodeServer(grpcServer, n)
	n.logger.Infow("node running on port", "port", listenAddr)

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
	fmt.Println("received tx from:", peer)

	return &pb.Ack{}, nil
}

func getPortNum(listenAddr string) (string, error) {
	values := strings.Split(listenAddr, ":")
	if len(values) != 2 {
		return "", errors.New("invalid listen value")
	}

	return values[1], nil
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
