package dcrypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
)

type PrivateKey struct {
	Key ed25519.PrivateKey
}

type PublicKey struct {
	Key ed25519.PublicKey
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

func (p *PrivateKey) Bytes() []byte {
	return p.Key
}

func (p *PrivateKey) Sign(msg []byte) []byte {
	return ed25519.Sign(p.Key, msg)
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, 32) // 32 bytes is length of public key
	copy(b, p.Key[32:])

	return &PublicKey{
		Key: b,
	}
}
