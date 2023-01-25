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
	// friendList  []string
	KnownConns  []string
	activeConns []*websocket.Conn // player
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
	for _, address := range s.KnownConns {
		s.ConnectOne(address)
	}
}

// == RECEIVE FUNCS ==
func (s *Server) HandleConn(ws *websocket.Conn) { // handleNewPlayer
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())
	status := checkConnectionPossible(ws, s)
	statusEx := checkConnectionExists(ws, s)

	// if less than 9 conns and conn dont exists
	// true + false
	if status && !statusEx {
		ws.Write(types.NewMessage(s.IP, []byte("Connecting...")))
		ac := NewConnectionAction("active", ws)
		s.ReceiveConnch(ac)
		s.readLoop(ws)

	} else {
		ws.Close()
	}
}

func (s *Server) ReceiveConnch(conna *ConnAction) {
	s.connch <- conna
}

func (s *Server) RunConnectionLoop() { // loop
	fmt.Println("Connection loop started ...")
	for conna := range s.connch {
		s.HandleConnectionAction(conna)
	}
}

func (s *Server) HandleConnectionAction(conna *ConnAction) { //handle message
	switch conna.action {
	case "active":
		s.addConnectionToServer(conna.conn)
	case "known":
		saveToKnown(s, conna.conn.RemoteAddr().String())
	case "delete":
		deleteConnection(s, conna.conn)
	default:
		panic("Invalid message")
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

				ac := NewConnectionAction("delete", ws)
				s.ReceiveConnch(ac)

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

		response := types.NewMessage(s.IP, []byte("Message received"))
		ws.Write(response)
	}
}

func (s *Server) addConnectionToServer(ws *websocket.Conn) { // add player
	s.activeConns = append(s.activeConns, ws)
	s.ConnCounter += 1

	// try to add to known
	ac := NewConnectionAction("known", ws)
	s.ReceiveConnch(ac)
}
