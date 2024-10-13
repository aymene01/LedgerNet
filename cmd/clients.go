package main

import (
	"github.com/aymene01/ledgerNet/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getClient(listenAddr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	client, err := grpc.NewClient(listenAddr, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getNodeClient(client *grpc.ClientConn) pb.NodeClient {
	return pb.NewNodeClient(client)
}
