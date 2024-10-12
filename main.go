package main

import (
	"time"

	"github.com/aymene01/ledgerNet/crypto"
	"github.com/aymene01/ledgerNet/node"
)

func main() {
	makeNode(":3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"}, false)
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"}, false)

	time.Sleep(time.Second)

	for {
		time.Sleep(time.Second * 2)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootsrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "blocker-1",
		ListenAddr: listenAddr,
	}

	if isValidator {
		cfg.PrivateKey = crypto.GeneratePrivateKey()
	}

	n := node.NewNode(cfg)
	go n.Start(listenAddr, bootsrapNodes)
	return n
}
