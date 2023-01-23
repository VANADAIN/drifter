package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

type Server struct {
	Name        string
	ID          dcrypto.PublicKey
	IP          string
	ConnCounter int
	// aliases     map[string]string // for local names
	friendList  []string
	knownConns  []string
	activeConns []*websocket.Conn
	connch      chan *ConnAction
}

func NewServer(name string) *Server {
	server := &Server{
		Name:        name,
		ConnCounter: 0,
		activeConns: make([]*websocket.Conn, 0),
		connch:      make(chan *ConnAction, 10),
	}

	go server.RunConnectionLoop()

	return server
}

// == DIAL FUNC ==

func (s *Server) ConnectOne(address string) {
	// connect to default endpoint
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		fmt.Println("Error in reading ip address")
	}

	ws_url := "ws" + strings.TrimPrefix(host, "http") + "/ws"
	ws, err := websocket.Dial(ws_url, "", host+":3000")
	if err != nil {
		fmt.Printf("Error conencting to: %s", host)
		return
	}

	s.readLoop(ws)
}

func (s *Server) CreateRandomConnections() {
	for _, address := range s.knownConns {
		s.ConnectOne(address)
	}
}

// == RECEIVE FUNCS ==
func (s *Server) HandleConn(ws *websocket.Conn) {
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())
	status := checkConnectionPossible(ws, s)
	statusEx := checkConnectionExists(ws, s)

	// if less than 9 conns and conn dont exists
	// true + false
	if status && !statusEx {
		ws.Write(types.NewMessage(s.IP, "Connecting..."))

		// add to active actions
		ac := NewConnection("active", ws)

		s.connch <- ac
		s.readLoop(ws)

	} else {
		ws.Close()
	}
}

func (s *Server) RunConnectionLoop() {
	fmt.Println("Connection loop started ...")
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
				ws.Close()
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}

		ms := types.Message{}
		msgraw := buf[:n]
		json.Unmarshal(msgraw, &ms)

		// decide what to do with this message ...
		// routes.Route(&ms)

		fmt.Println(ms.Body.Payload)

		response := types.NewMessage(s.IP, "Message received")
		ws.Write(response)
	}
}

func (s *Server) addConnection(ws *websocket.Conn) {
	s.activeConns = append(s.activeConns, ws)
	s.ConnCounter += 1

	// if connection was added to active after all checks
	// try to add to known
	// send to ch with action known
	ac := NewConnection("known", ws)

	s.connch <- ac
}
