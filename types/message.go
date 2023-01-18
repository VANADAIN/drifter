package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Message struct {
	Header MsgHeader
	Body   MsgBody
}

type MsgBody struct {
	Type    string `json: "type"`
	From    string `json: "from"`
	Origin  string `json: "origin"`
	Payload string `json: "payload"`
}

type MsgHeader struct {
	Hash      [32]byte
	CreatedAt time.Time
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
