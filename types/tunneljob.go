package types

import (
	"encoding/json"
	"time"
)

type TunnelJob Message

func NewTJ(from string, payload []byte) []byte {
	msg := &Message{
		Header: MsgHeader{
			CreatedAt: time.Now().Unix(),
		},
		Body: MsgBody{
			Type:    "tunneljob",
			From:    from,
			Origin:  from,
			Payload: payload,
		},
	}

	msgj, _ := json.Marshal(msg)

	return msgj
}
