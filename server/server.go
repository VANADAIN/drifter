package server

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/net/websocket"
)

type Server struct {
	known       []*net.Addr
	activeConns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		activeConns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleConn(ws *websocket.Conn) {
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())

	// lock actor
	s.activeConns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remote connection closed
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		ws.Write([]byte("msg received"))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.activeConns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}
