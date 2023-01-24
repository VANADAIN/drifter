package keys_storage

import (
	"fmt"
	"os"

	"github.com/VANADAIN/drifter/dcrypto"
)

type KeysStorage struct {
	KeysPath string
}

func NewKeysStorage(path string) *KeysStorage {
	store := &KeysStorage{
		KeysPath: path,
	}

	return store
}

func (ks *KeysStorage) SaveKeys(pk *dcrypto.PrivateKey, pubk *dcrypto.PublicKey) {
	err := os.WriteFile(ks.KeysPath+"pk", pk.Bytes(), 0666)
	if err != nil {
		panic("Error saving crypto keys")
	}

	err2 := os.WriteFile(ks.KeysPath+"pubk", pubk.Bytes(), 0666)
	if err2 != nil {
		panic("Error saving crypto keys")
	}
}

func (ks *KeysStorage) LoadKeys() (*dcrypto.PrivateKey, *dcrypto.PublicKey) {
	// check if file exists
	pk_b, err := os.ReadFile(ks.KeysPath + "pk")
	if err != nil {
		if os.IsNotExist(err) {
			panic("Error loading crypto keys.")
		} else {
			panic(fmt.Sprintf("Error reading crypto keys %s:", err))
		}
	}

	pk := &dcrypto.PrivateKey{
		Key: pk_b,
	}
	pub := pk.Public()

	return pk, pub
}
