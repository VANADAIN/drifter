package server

import "golang.org/x/net/websocket"

func checkConnectionPossible(ws *websocket.Conn, s *Server) bool {
	// max of connection active
	if len(s.activeConns) < 10 {
		return true
	}

	ws.Write([]byte("This node reached maximum number of connections. Closing connection..."))
	return false
}

func checkConnectionExists(ws *websocket.Conn, s *Server) bool {
	if s.activeAddr[ws.RemoteAddr().String()] {
		return true
	}

	ws.Write([]byte("Your node is already connected. Closing connection..."))
	return false
}
