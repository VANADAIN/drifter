package node

import (
	"crypto/ed25519"
	"encoding/binary"
	"log"
	"os"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// save private key to device
func saveKeyToFile(key *ed25519.PrivateKey) {
	file, err := os.Create("./keys/key.private")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, key)
	if err != nil {
		log.Fatal("Private key saving failed")
	}
}

func readKeyFromFile() ed25519.PrivateKey {
	f, err := os.Open("./keys/key.private")
	if err != nil {
		log.Fatal("Private key file not found")
	}

	defer f.Close()

	// var key ed25519.PrivateKey
	key := make([]byte, 64)
	err = binary.Read(f, binary.LittleEndian, &key)
	if err != nil {
		log.Fatal(err)
	}

	return key
}
