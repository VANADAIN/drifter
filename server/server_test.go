package server

import (
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestConns(t *testing.T) {
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
	defer server.Close()

	fmt.Println("Server url: ", server.URL)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	origin := "http://127.0.0.1:50000"
	conn, err := websocket.Dial(wsURL, "", origin)

	if err != nil {
		panic(err)
	}

	fmt.Println("Local address of client:", conn.LocalAddr().String())

	time.Sleep(2 * time.Second)

	fmt.Printf("\nServer: %+v\n\n", server)
	fmt.Printf("\nConnection: %+v\n\n", conn)

	assert.Equal(t, node_serv.KnownConns[0], origin)
	assert.Equal(t, node_serv.activeConns[0].RemoteAddr().String(), origin)
	assert.Equal(t, node_serv.ConnCounter, 1)
}
