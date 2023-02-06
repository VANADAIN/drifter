package types

import (
	"encoding/json"
	"time"
)

type Barter Message

func NewBarter(from string, payload []byte) []byte {
	msg := &Message{
		Header: MsgHeader{
			CreatedAt: time.Now().Unix(),
		},
		Body: MsgBody{
			Type:    "Barter",
			From:    from,
			Origin:  from,
			Payload: payload,
		},
	}

	msgj, _ := json.Marshal(msg)

	return msgj
}
