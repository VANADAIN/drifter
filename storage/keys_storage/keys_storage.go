package keys_storage

import (
	"fmt"
	"os"
	"path/filepath"

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
	ks.CreatePath()
	err := os.WriteFile(ks.KeysPath+"pk", pk.Bytes(), 0666)
	// fmt.Println(err)
	if err != nil {
		panic("Error saving private key")
	}

	err2 := os.WriteFile(ks.KeysPath+"pubk", pubk.Bytes(), 0666)
	if err2 != nil {
		panic("Error saving public key")
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

func (ks *KeysStorage) CreatePath() {
	// create path
	newpath := filepath.Join(".", ks.KeysPath)
	e := os.MkdirAll(newpath, os.ModePerm)
	if e != nil {
		panic(e)
	}
}
