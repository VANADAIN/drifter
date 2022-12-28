package network

import (
	"net"

	"github.com/ethereum/go-ethereum/log"
)

type Mediator struct {
	nm *NetworkManager // to access data
	c  net.Conn
}

func (m *Mediator) RunReadLoop() {
	//
}

func (m *Mediator) SendNodeInfo() {
	nodeinfo := m.nm.info.Bytes()
	_, err := m.c.Write(nodeinfo)
	if err != nil {
		log.Errorf("Unable to send node info to: %s", m.c.RemoteAddr().String())
	}
}

func (m *Mediator) SendMessage() {
	// create message
	// ...
	// m.c.Write(msg)
}
