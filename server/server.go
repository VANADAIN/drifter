package server

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

type Server struct {
	Name        string
	ID          dcrypto.PublicKey
	connCounter int
	knownConns  []string
	activeConns map[*websocket.Conn]bool
	activeAddr  map[string]bool
	connch      chan *ConnAction
}

type ConnAction struct {
	action string
	conn   *websocket.Conn
}

func NewServer() *Server {
	server := &Server{
		activeConns: make(map[*websocket.Conn]bool),
		activeAddr:  make(map[string]bool),
		connch:      make(chan *ConnAction, 10),
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

		// add to active actions
		ac := &ConnAction{
			action: "active",
			conn:   ws,
		}

		s.connch <- ac
		s.readLoop(ws)

	} else {
		ws.Close()
	}
}

func (s *Server) RunConnectionLoop() {
	for conna := range s.connch {
		switch conna.action {
		case "active":
			s.addConnection(conna.conn)
		case "known":
			go saveToKnown(s, conna.conn.RemoteAddr().String())
		}
	}
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

		ms := types.Message{}
		msgraw := buf[:n]
		json.Unmarshal(msgraw, &ms)

		// msg := types.Message{}
		// websocket.JSON.Receive(ws, &msg)

		//routes.Route(&msg)

		fmt.Println("msg from req")
		fmt.Println(ms.Body.Payload)

		respp := types.Message{
			Header: types.MsgHeader{
				CreatedAt: time.Now().Unix(),
			},
			Body: types.MsgBody{
				Type:    "text",
				Payload: "Msg received",
			},
		}
		resp, _ := json.Marshal(respp)
		ws.Write(resp)
	}
}

func (s *Server) addConnection(conn *websocket.Conn) {
	s.activeConns[conn] = true
	s.activeAddr[conn.RemoteAddr().String()] = true
	s.connCounter += 1

	// if connection was added to active after all checks
	// try to add to known
	// send to ch with action known
	ac := &ConnAction{
		action: "known",
		conn:   conn,
	}

	s.connch <- ac
}
