package types

import (
	"bytes"
	"encoding/gob"

	"github.com/VANADAIN/drifter/dcrypto"
)

type NodeInfo struct {
	ID   *dcrypto.PublicKey
	Desc string
}

func (ni *NodeInfo) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(ni)

	return buf.Bytes()
}
