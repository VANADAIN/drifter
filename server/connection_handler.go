package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/VANADAIN/drifter/types"
	"golang.org/x/net/websocket"
)

type ConnectionHandler struct {
	s           *Server
	ConnCounter int
	// aliases     map[string]string // for local names
	// friendList  []string or ws ???
	KnownConns  []string
	activeConns []*websocket.Conn
	connch      chan *ConnAction
}

// == DIAL FUNC ==

func (ch *ConnectionHandler) ConnectOne(address string) {
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

	ch.readLoop(ws)
}

func (ch *ConnectionHandler) CreateRandomConnections() {
	for _, address := range ch.KnownConns {
		ch.ConnectOne(address)
	}
}

// == RECEIVE FUNCS ==
func (ch *ConnectionHandler) HandleConn(ws *websocket.Conn) {
	fmt.Println("New incoming conn from: ", ws.RemoteAddr())
	status := checkConnectionPossible(ws, ch)
	statusEx := checkConnectionExists(ws, ch)

	// if less than 9 conns and conn dont exists
	// true + false
	if status && !statusEx {
		ws.Write(types.NewMessage(ch.s.IP, []byte("Connecting...")))
		ac := NewConnectionAction("active", ws)
		ch.ReceiveConnch(ac)
		ch.readLoop(ws)

	} else {
		ws.Close()
	}
}

func (ch *ConnectionHandler) ReceiveConnch(conna *ConnAction) {
	ch.connch <- conna
}

func (ch *ConnectionHandler) RunConnectionLoop() {
	fmt.Println("Connection loop started ...")
	for conna := range ch.connch {
		ch.HandleConnectionAction(conna)
	}
}

func (ch *ConnectionHandler) HandleConnectionAction(conna *ConnAction) {
	switch conna.action {
	case "active":
		ch.addConnectionToServer(conna.conn)
	case "known":
		saveToKnown(ch, conna.conn.RemoteAddr().String())
	case "delete":
		deleteConnection(ch, conna.conn)
	default:
		panic("Invalid message")
	}
}

func (ch *ConnectionHandler) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remote connection closed
				ws.Close()
				ac := NewConnectionAction("delete", ws)
				ch.ReceiveConnch(ac)

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

		response := types.NewMessage(ch.s.IP, []byte("Message received"))
		ws.Write(response)
	}
}

func (ch *ConnectionHandler) addConnectionToServer(ws *websocket.Conn) {
	ch.activeConns = append(ch.activeConns, ws)
	ch.ConnCounter += 1

	// try to add to known
	ac := NewConnectionAction("known", ws)
	ch.ReceiveConnch(ac)
}
