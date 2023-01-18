package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestConns(t *testing.T) {
	server := NewServer()

	handler := websocket.Handler(server.HandleConn)
	serv := httptest.NewServer(http.Handler(handler))

	fmt.Println(serv.URL)

	address := serv.URL
	wsURL := "ws" + strings.TrimPrefix(serv.URL, "http") + "/ws"

	_, err := websocket.Dial(wsURL, "", address)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%+v\n", server)
	// fmt.Printf("%+v\n", conn)

	time.Sleep(2 * time.Second)

	assert.Equal(t, server.connCounter, 1)
	assert.Equal(t, server.activeAddr[serv.URL], true)
	assert.Equal(t, server.knownConns[0], serv.URL)
	// assert.Equal(t, server.activeConns[conn], true) ???
}
