package main

import (
	"log"
	"time"

	"github.com/aymene01/blocker/node"
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
