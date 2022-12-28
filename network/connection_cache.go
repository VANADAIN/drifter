package network

import "github.com/VANADAIN/drifter/types"

type ConnCache struct {
	history []*types.Message
}

func (c *ConnCache) Add(m *types.Message) {
	c.history = append(c.history, m)
}

func (c *ConnCache) Flush() {
	c.history = nil
}
