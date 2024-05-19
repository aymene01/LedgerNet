package main

import (
	"context"
	"log"

	"github.com/aymene01/blocker/pb"
)

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
