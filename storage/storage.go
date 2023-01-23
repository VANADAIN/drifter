package storage

import (
	"fmt"
	"os"

	"github.com/VANADAIN/drifter/dcrypto"
)

type StorageManager struct {
	keysPath       string
	connectionPath string
}

func (sm *StorageManager) SaveKeys(pk *dcrypto.PrivateKey, pubk *dcrypto.PublicKey) {
	err := os.WriteFile("data/pk", pk.Bytes(), 0666)
	if err != nil {
		panic("Error saving crypto keys")
	}

	err2 := os.WriteFile("data/pubk", pubk.Bytes(), 0666)
	if err2 != nil {
		panic("Error saving crypto keys")
	}
}

func (sm *StorageManager) LoadKeys() (*dcrypto.PrivateKey, *dcrypto.PublicKey) {
	// check if file exists
	pkey_b, err := os.ReadFile("data/pk")
	if err != nil {
		if os.IsNotExist(err) {
			panic("Error loading crypto keys.")
		} else {
			panic(fmt.Sprintf("Error reading crypto keys %s:", err))
		}
	}

	pk := &dcrypto.PrivateKey{
		Key: pkey_b,
	}
	pub := pk.Public()

	return pk, pub
}
