package main

import (
	"time"

	"github.com/aymene01/ledgerNet/node"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})

	time.Sleep(4 * time.Second)
	makeNode(":5000", []string{":4000"})

	select {}
}

func makeNode(listenAddr string, bootsrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr, bootsrapNodes)
	return n
}
