package cache

import "github.com/VANADAIN/drifter/types"

type MessageCache struct {
	Data []*types.Message
}

func NewMessageCache() *MessageCache {
	return &MessageCache{
		Data: make([]*types.Message, 0),
	}
}
