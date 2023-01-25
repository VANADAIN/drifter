package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func CreateServer() (*httptest.Server, *Server) {
	node_serv := NewServer("A")
	handler := websocket.Handler(node_serv.HandleConn)
	server := httptest.NewUnstartedServer(handler)

	l, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}

	server.Listener.Close()
	server.Listener = l

	server.Start()

	return server, node_serv
}

func TestConnections(t *testing.T) {
	server, node_serv := CreateServer()
	defer server.Close()

	fmt.Println("Server url: ", server.URL)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	origin := "http://127.0.0.1:50000"
	conn, err := websocket.Dial(wsURL, "", origin)

	if err != nil {
		panic(err)
	}

	fmt.Println("Local address of client:", conn.LocalAddr().String())

	time.Sleep(1 * time.Second)

	fmt.Printf("\nServer: %+v\n\n", server)
	fmt.Printf("\nConnection: %+v\n\n", conn)

	assert.Equal(t, node_serv.KnownConns[0], origin)
	assert.Equal(t, node_serv.activeConns[0].RemoteAddr().String(), origin)
	assert.Equal(t, node_serv.ConnCounter, 1)
}

func TestConnectionDenial(t *testing.T) {
	server, _ := CreateServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	origin := "http://127.0.0.1:50000"
	_, err := websocket.Dial(wsURL, "", origin)

	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	conn, err := websocket.Dial(wsURL, "", origin)

	buf := make([]byte, 1024)
	var check error
	for {
		_, err := conn.Read(buf)
		check = err
		if err != nil {
			break
		}
	}

	assert.NotNil(t, check)
	assert.Equal(t, check, io.EOF)
}
