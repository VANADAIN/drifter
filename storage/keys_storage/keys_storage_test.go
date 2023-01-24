package keys_storage

import (
	"os"
	"testing"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/stretchr/testify/assert"
)

func TestKeyStorage(t *testing.T) {
	store := NewKeysStorage("data/keys/")

	pk := dcrypto.GeneratePrivateKey()
	pub := pk.Public()

	store.SaveKeys(pk, pub)
	_, err := os.Stat("data/keys/pk")
	_, err2 := os.Stat("data/keys/pubk")

	assert.Nil(t, err)
	assert.Nil(t, err2)

	os.RemoveAll("data")
}
