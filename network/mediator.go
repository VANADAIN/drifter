package network

import (
	"fmt"
	"io"
	"net"

	"github.com/VANADAIN/drifter/types"
	"github.com/ethereum/go-ethereum/log"
)

type Mediator struct {
	nm *NetworkManager // to access data of parent nm
	c  net.Conn
}

func (m *Mediator) RunReadLoop() {
	defer m.c.Close()

	buf := make([]byte, 2048)
	for {
		n, err := m.c.Read(buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("Connection closed")
				m.nm.DeleteActiveConnection(m.c)
				break
			} else {
				fmt.Println("Read error: ", err)
			}
		}

		msgb := types.MessageBody{
			From:    m.c.RemoteAddr().String(),
			Payload: buf[:n],
		}

		msg := types.Message{
			Body: msgb,
		}
		msg.HashIt()

		// fmt.Println(msg)

		m.nm.Msgch <- msg

		m.c.Write([]byte("msg received"))
	}
}

// ========== service functions ==================

func (m *Mediator) SendNodeInfo() {
	nodeinfo := m.nm.info.Bytes()
	_, err := m.c.Write(nodeinfo)
	if err != nil {
		log.Errorf("Unable to send node info to: %s", m.c.RemoteAddr().String())
	}
}

func (m *Mediator) RequestPublicKey() {
	msg := &types.Message{
		Header: types.MessageHeader{
			Type: "Public key request",
		},
		Body: types.MessageBody{
			From: m.c.LocalAddr().String(),
		},
	}
	b := msg.Bytes()

	m.c.Write(b)
}

// ============ default functions ================

func (m *Mediator) SendMessage() {
	// create message
	// ...
	// m.c.Write(msg)
}
