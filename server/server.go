package server

import (
	"fmt"
	"io"

	"github.com/VANADAIN/drifter/routes"
	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

type Server struct {
	connCounter int
	knownConns  []string
	activeConns map[*websocket.Conn]bool
	activeAddr  map[string]bool
	connch      chan *websocket.Conn
}

func NewServer() *Server {
	server := &Server{
		activeConns: make(map[*websocket.Conn]bool),
		activeAddr:  make(map[string]bool),
		connch:      make(chan *websocket.Conn, 10),
	}

	go server.RunConnectionLoop()

	return server
}

func (s *Server) HandleConn(ws *websocket.Conn) {
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())
	status := checkConnectionPossible(ws, s)
	statusEx := checkConnectionExists(ws, s)

	// if less than 9 conns and conn dont exists
	// true + false
	if status && !statusEx {
		ws.Write([]byte("Connecting..."))
		s.connch <- ws
		s.readLoop(ws)
	} else {
		ws.Close()
	}
}

func (s *Server) RunConnectionLoop() {
	for conn := range s.connch {
		s.addConnection(conn)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remote connection closed
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}

		msg := types.Message{}
		websocket.JSON.Receive(ws, &msg)
		routes.Route(&msg)

		fmt.Println(msg.Body.Payload)

		ws.Write([]byte(fmt.Sprintf("msg received %d", s.connCounter)))
	}
}

func (s *Server) broadcast(msg *types.Message) {
	for ws := range s.activeConns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg.Bytes()); err != nil {
				fmt.Println("Broadcast error: ", err)
			}
		}(ws)
	}
}

func (s *Server) addConnection(conn *websocket.Conn) {
	s.activeConns[conn] = true
	s.activeAddr[conn.RemoteAddr().String()] = true
	s.connCounter += 1
}
