package node

import (
	"context"
	"fmt"

	"github.com/aymene01/blocker/pb"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	pb.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: "blocker-1",
	}
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
