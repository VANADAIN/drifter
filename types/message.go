package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Message struct {
	Header MsgHeader `json: "header"`
	Body   MsgBody   `json: "body"`
}

type MsgBody struct {
	Type    string `json: "type"`
	From    string `json: "from"`
	Origin  string `json: "origin"`
	Payload string `json: "payload"`
}

type MsgHeader struct {
	Hash      [32]byte `json: "hash"`
	CreatedAt int64    `json: "createdAt"`
}

func SimpleTextMessage(from string, payload string) *Message {
	msg := &Message{
		Header: MsgHeader{
			CreatedAt: time.Now().Unix(),
		},
		Body: MsgBody{
			Type:    "text",
			From:    from,
			Origin:  from,
			Payload: payload,
		},
	}

	return msg
}

func (m *Message) HashIt() {
	b := m.Body.Bytes()
	h := sha256.Sum256(b)
	m.Header.Hash = h
}

func (m *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(m)

	return buf.Bytes()
}

func (m *MsgBody) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(m)

	return buf.Bytes()
}
