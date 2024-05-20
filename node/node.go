package node

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/aymene01/ledgerNet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string

	peerLock sync.RWMutex
	peers    map[pb.NodeClient]bool

	pb.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		peers:   make(map[pb.NodeClient]bool),
		version: "blocker-1",
	}
}

func (n *Node) addPeer(c pb.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.peers[c] = true
}

func (n *Node) deletePeer(c pb.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, c)
}

func (n *Node) Start(listenAddr string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return err
	}

	pb.RegisterNodeServer(grpcServer, n)

	portNum, err := getPortNum(listenAddr)

	if err != nil {
		return err
	}

	fmt.Println("node running on port:", portNum)

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *pb.Version) (*pb.Version, error) {
	ourVersion := &pb.Version{
		Version: n.version,
		Height:  100,
	}

	p, _ := peer.FromContext(ctx)
	fmt.Printf("received version from %s: %+v\n", v, p.Addr)

	return ourVersion, nil
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
