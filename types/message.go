package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// this is the file for one-time interaction methods
type Message struct {
	DataHash [32]byte
	Body     MessageBody
}

type MessageBody struct {
	From    string
	Payload []byte
}

func (m *Message) AddHash() {
	b := m.Body.Bytes()
	h := sha256.Sum256(b)
	m.DataHash = h
}

func (body *MessageBody) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(body)

	return buf.Bytes()
}
