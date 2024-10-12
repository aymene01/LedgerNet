package main

import (
	"context"
	"log"

	"github.com/aymene01/ledgerNet/crypto"

	"github.com/aymene01/ledgerNet/pb"
	"github.com/aymene01/ledgerNet/util"
)

const currentListenAddr = ":3000"

func makeHandshake() {
	client, err := getClient(currentListenAddr)

	if err != nil {
		log.Fatal(err)
	}

	c := getNodeClient(client)
	version := &pb.Version{
		Version:    "Blocker-0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), version)

	if err != nil {
		log.Fatal(err)
	}
}

func makeTransaction() {
	client, err := getClient(currentListenAddr)

	if err != nil {
		log.Fatal(err)
	}

	c := getNodeClient(client)

	privateKey := crypto.GeneratePrivateKey()
	
	tx := &pb.Transaction{
		Version: 1,
		Inputs: []*pb.TxInput{
			{
				PrevHash: util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey: privateKey.Public().Bytes(),
			},
		},
		Outputs: []*pb.TxOutput{
			{
				Amount: 99,
				Address: privateKey.Public().Address().Bytes(),
			},
		},
	}
	
	_, err = c.HandleTransaction(context.TODO(), tx)

	if err != nil {
		log.Fatal(err)
	}
}
