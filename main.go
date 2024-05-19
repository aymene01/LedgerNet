package main

import (
	"context"
	"log"
	"time"

	"github.com/aymene01/blocker/node"
	"github.com/aymene01/blocker/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	node := node.NewNode()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeHandshake()
		}
	}()

	log.Fatal(node.Start(":3000"))
}

func getClient() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	client, err := grpc.NewClient(":3000", opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getNodeClient(client *grpc.ClientConn) pb.NodeClient {
	return pb.NewNodeClient(client)
}

func makeHandshake() {
	client, err := getClient()

	if err != nil {
		log.Fatal(err)
	}

	c := getNodeClient(client)
	version := &pb.Version{
		Version: "Blocker-0.1",
		Height:  1,
	}

	_, err = c.Handshake(context.TODO(), version)

	if err != nil {
		log.Fatal(err)
	}
}

func makeTransaction() {
	client, err := getClient()

	if err != nil {
		log.Fatal(err)
	}

	c := getNodeClient(client)
	_, err = c.HandleTransaction(context.TODO(), &pb.Transaction{})

	if err != nil {
		log.Fatal(err)
	}
}
