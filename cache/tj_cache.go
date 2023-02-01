package cache

import "github.com/VANADAIN/drifter/types"

type TJCache struct {
	Data []*types.Message
}

func NewTJCache() *TJCache {
	return &TJCache{
		Data: make([]*types.Message, 0),
	}
}
