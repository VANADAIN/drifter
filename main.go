package main

import (
	"net/http"

	"github.com/VANADAIN/drifter/server"
	"golang.org/x/net/websocket"
)

func main() {

	// ctx := context.Background()
	// err := network.GetIPAndForwardPort(ctx)
	// if err == nil {
	// 	fmt.Println("port forwarding created!")
	// }

	server := server.NewServer()
	http.Handle("/ws", websocket.Handler(server.HandleConn))
	http.ListenAndServe(":3000", nil)
}
