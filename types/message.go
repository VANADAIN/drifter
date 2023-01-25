package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
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
	Payload []byte `json: "payload"`
}

type MsgHeader struct {
	Hash      [32]byte `json: "hash"`
	CreatedAt int64    `json: "createdAt"`
}

func NewMessage(from string, payload []byte) []byte {
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

	msgj, _ := json.Marshal(msg)

	return msgj
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
