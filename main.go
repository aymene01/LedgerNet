package main

import (
	"log"

	"github.com/aymene01/ledgerNet/node"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		makeHandshake()
	// 	}
	// }()
	select {}
}

func makeNode(listenAddr string, bootsrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr)
	if len(bootsrapNodes) > 0 {
		if err := n.BootstrapNetwork(bootsrapNodes); err != nil {
			log.Fatal(err)
		}
	}

	return n
}
