package server

import (
	"fmt"

	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

func Send(s *Server, msg *types.Message) {

}

func Broadcast(s *Server, msg *types.Message) {
	for ws := range s.activeConns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg.Bytes()); err != nil {
				fmt.Println("Broadcast error: ", err)
			}
		}(ws)
	}
}
