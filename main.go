package main

import (
	"log"
	"time"

	"github.com/aymene01/ledgerNet/node"
)

func main() {
	node := node.NewNode()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeTransaction()
		}
	}()

	log.Fatal(node.Start(":3000"))
}
