package server

import (
	"fmt"

	"encoding/json"

	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

func Send(s *Server, msg *types.Message) {

}

func Broadcast(s *Server, msg *types.Message) {
	for _, ws := range s.CH.activeConns {
		go func(ws *websocket.Conn) {
			msgj, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Json error: ", err)
			}

			if _, err := ws.Write(msgj); err != nil {
				fmt.Println("Broadcast error: ", err)
			}
		}(ws)
	}
}
