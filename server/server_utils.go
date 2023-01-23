package server

import (
	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

func checkConnectionPossible(ws *websocket.Conn, s *Server) bool {
	// max of connection active
	if len(s.activeConns) < 10 {
		return true
	}

	msg := types.NewMessage(s.IP, "This node reached maximum number of connections. Closing connection...")
	ws.Write(msg)

	return false
}

func checkConnectionExists(ws *websocket.Conn, s *Server) bool {
	for _, conn := range s.activeConns {
		if ws == conn {
			payload := "Your node is already connected. Closing connection..."
			msg := types.NewMessage(s.IP, payload)
			ws.Write(msg)

			return true
		}
	}

	return false
}

func saveToKnown(s *Server, addr string) {
	saved := addrSaved(s, addr)
	if !saved {
		saveAddr(s, addr)
	} else {
		return
	}
}

func saveAddr(s *Server, addr string) {
	// concurrent safe ???
	s.knownConns = append(s.knownConns, addr)
}

func addrSaved(s *Server, addr string) bool {
	for _, val := range s.knownConns {
		if val == addr {
			return true
		}
	}
	return false
}
