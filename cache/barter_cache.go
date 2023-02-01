package cache

import "github.com/VANADAIN/drifter/types"

type BarterCache struct {
	Data []*types.Message
}

func NewBarterCache() *BarterCache {
	return &BarterCache{
		Data: make([]*types.Message, 0),
	}
}
