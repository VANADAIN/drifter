package server

import (
	"fmt"

	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

func checkConnectionPossible(ws *websocket.Conn, ch *ConnectionHandler) bool {
	// max of connection active
	if len(ch.activeConns) < 10 {
		return true
	}

	msg := types.NewMessage(ch.s.IP, []byte("This node reached maximum number of connections. Closing connection..."))
	ws.Write(msg)

	return false
}

func checkConnectionExists(ws *websocket.Conn, ch *ConnectionHandler) bool {
	for _, conn := range ch.activeConns {
		if ws.RemoteAddr().String() == conn.RemoteAddr().String() {
			fmt.Println("Connection already exists")
			payload := "Your node is already connected. Closing connection..."
			msg := types.NewMessage(ch.s.IP, []byte(payload))
			ws.Write(msg)

			return true
		}
	}

	return false
}

func saveToKnown(s *ConnectionHandler, addr string) {
	saved := addrSaved(s, addr)
	if !saved {
		saveAddr(s, addr)
	} else {
		return
	}
}

func saveAddr(s *ConnectionHandler, addr string) {
	// concurrent safe ???
	s.KnownConns = append(s.KnownConns, addr)
}

func addrSaved(s *ConnectionHandler, addr string) bool {
	for _, val := range s.KnownConns {
		if val == addr {
			return true
		}
	}
	return false
}

func deleteConnection(s *ConnectionHandler, ws *websocket.Conn) {
	for index, conn := range s.activeConns {
		if conn.RemoteAddr().String() == ws.RemoteAddr().String() {
			s.activeConns = append(s.activeConns[:index], s.activeConns[index+1:]...)
		}
	}
}
