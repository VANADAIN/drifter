package server

import (
	"golang.org/x/net/websocket"
)

func checkConnectionPossible(ws *websocket.Conn, s *Server) bool {
	// max of connection active
	if len(s.activeConns) < 10 {
		return true
	}

	ws.Write([]byte("This node reached maximum number of connections. Closing connection..."))
	return false
}

func checkConnectionExists(ws *websocket.Conn, s *Server) bool {
	if !s.activeAddr[ws.RemoteAddr().String()] {
		return false
	}

	ws.Write([]byte("Your node is already connected. Closing connection..."))
	return true
}

func saveToKnown(s *Server, addr string) {
	saved := addrSaved(s, addr)
	if !saved {
		// not in list
		saveAddr(s, addr)
	} else {
		return
	}
}

func saveAddr(s *Server, addr string) {
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
