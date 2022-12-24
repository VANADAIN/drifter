package dcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), 64)

	pub := privKey.Public()
	assert.Equal(t, len(pub.Bytes()), 32)
}

func TestPrivateKeySign(t *testing.T) {
	priv := GeneratePrivateKey()
	pub := priv.Public()

	msg := []byte("hello")
	sig := priv.Sign(msg)
	assert.True(t, sig.Verify(pub, msg))

	assert.False(t, sig.Verify(pub, []byte("foo")))

	priv2 := GeneratePrivateKey()
	pub2 := priv2.Public()
	assert.False(t, sig.Verify(pub2, msg))
}
