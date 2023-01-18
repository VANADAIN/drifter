package server

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type Server struct {
	connCounter int
	activeConns map[*websocket.Conn]bool
	connch      chan *websocket.Conn
}

func NewServer() *Server {
	server := &Server{
		activeConns: make(map[*websocket.Conn]bool),
		connch:      make(chan *websocket.Conn, 10),
	}

	go server.RunConnectionLoop()

	return server
}

func (s *Server) HandleConn(ws *websocket.Conn) {
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())
	s.connch <- ws
	s.readLoop(ws)
}

func (s *Server) RunConnectionLoop() {
	for conn := range s.connch {
		s.addConnection(conn)
	}
}

func (s *Server) addConnection(conn *websocket.Conn) {
	fmt.Println("get conn from connch ", conn.RemoteAddr())
	s.activeConns[conn] = true
	s.connCounter += 1
}

func (s *Server) readLoop(ws *websocket.Conn) {
	fmt.Println("Reading from: ", ws.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remote connection closed
				break
			}
			fmt.Println("Read error: ", err)
			break
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		ws.Write([]byte(fmt.Sprintf("msg received %d", s.connCounter)))
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
