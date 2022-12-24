package dcrypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type PrivateKey struct {
	Key ed25519.PrivateKey
}

type PublicKey struct {
	Key ed25519.PublicKey
}

type Signature struct {
	value []byte
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		Key: ed25519.NewKeyFromSeed(seed),
	}
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != 32 {
		panic("Keys: invalid seed lenght")
	}

	return &PrivateKey{
		Key: ed25519.NewKeyFromSeed(seed),
	}
}

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func (p *PrivateKey) Bytes() []byte {
	return p.Key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.Key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, 32) // 32 bytes is length of public key
	copy(b, p.Key[32:])

	return &PublicKey{
		Key: b,
	}
}

func (p *PrivateKey) String() string {
	return string(p.Bytes())
}

func (p *PublicKey) Bytes() []byte {
	return p.Key
}

func (p *PublicKey) String() string {
	return string(p.Bytes())
}

func (s *Signature) Verify(pub *PublicKey, msg []byte) bool {
	return ed25519.Verify(pub.Key, msg, s.value)
}
