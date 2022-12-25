package main

import (
	"context"
	"fmt"
	"log"

	"github.com/VANADAIN/drifter/network"
	"github.com/VANADAIN/drifter/node"
)

func main() {
	ctx := context.Background()
	err := network.GetIPAndForwardPort(ctx)
	if err == nil {
		fmt.Println("port forwarding created!")
	}

	server := node.NewNode(":3000")

	go func() {
		for msg := range server.Msgch {
			fmt.Printf("receiver msg from (%s): %s\n", msg.From, string(msg.Payload))
		}
	}()

	log.Fatal(network.Start(server))
}
