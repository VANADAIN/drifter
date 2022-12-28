package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// this is the file for peer 2 peer straight conversation
type Message struct {
	Header MessageHeader
	Body   MessageBody
}

type MessageHeader struct {
	Type     string
	DataHash [32]byte
}

type MessageBody struct {
	From    string
	Time    time.Time
	Payload []byte
}

func (m *Message) HashIt() {
	b := m.Body.Bytes()
	h := sha256.Sum256(b)
	m.Header.DataHash = h
}

func (m *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(m)

	return buf.Bytes()
}

func (body *MessageBody) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(body)

	return buf.Bytes()
}
