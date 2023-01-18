package types

type Message struct {
	Type    string `json: "type"`
	Payload string `json: "payload"`
}

// func (m *Message) HashIt() {
// 	b := m.payload.Bytes()
// 	h := sha256.Sum256(b)
// 	m.DataHash = h
// }

// func (m *Message) Bytes() []byte {
// 	buf := &bytes.Buffer{}
// 	enc := gob.NewEncoder(buf)
// 	enc.Encode(m)

// 	return buf.Bytes()
// }

// func (body *MessageBody) Bytes() []byte {
// 	buf := &bytes.Buffer{}
// 	enc := gob.NewEncoder(buf)
// 	enc.Encode(body)

// 	return buf.Bytes()
// }
