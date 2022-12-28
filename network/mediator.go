package network

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/VANADAIN/drifter/types"
)

type Mediator struct {
	nm    *NetworkManager // to access data of parent nm
	c     net.Conn
	cache *ConnCache
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

		// fmt.Println(msg)

		m.nm.Msgch <- msg

		m.c.Write([]byte("msg received"))
	}
}

// ========== service functions ==================

func (m *Mediator) SendNodeInfo() {
	nodeinfo := m.nm.info.Bytes()
	msg := &types.Message{
		Header: types.MessageHeader{
			Type: "Node Info",
		},
		Body: types.MessageBody{
			From:    m.c.LocalAddr().String(),
			Payload: nodeinfo,
		},
	}

	_, err := m.c.Write(msg.Bytes())
	if err != nil {
		log.Printf("Unable to send node info to: %s", m.c.RemoteAddr().String())
	}
}

func (m *Mediator) RegisterMe() {
	msg := &types.Message{
		Header: types.MessageHeader{
			Type: "Register Me",
		},
		Body: types.MessageBody{
			From: m.c.LocalAddr().String(),
		},
	}

	_, err := m.c.Write(msg.Bytes())
	if err != nil {
		log.Printf("Unable to send node info to: %s", m.c.RemoteAddr().String())
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

func (m *Mediator) SendMessage(message string) {
	msg := &types.Message{
		Header: types.MessageHeader{
			Type: "Default",
		},
		Body: types.MessageBody{
			From:    m.c.LocalAddr().String(),
			Payload: []byte(message),
		},
	}

	_, err := m.c.Write(msg.Bytes())
	if err != nil {
		log.Printf("Unable to send node info to: %s", m.c.RemoteAddr().String())
	}
}
