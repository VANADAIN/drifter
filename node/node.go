package node

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"log"
	"net"
	"os"
)

type Node struct {
	id         *ed25519.PublicKey
	listenPort string
	lsn        net.Listener
	msgch      chan Message
	quitch     chan struct{}
}

func NewNode(port string) *Node {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	status, err := exists("./keys")
	if err != nil {
		log.Fatal(err)
	}

	if status == false {
		err := os.Mkdir("./keys", 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		panic("Wrong entrypoint for node! Switch to Load mode.")
	}

	saveKeyToFile(&priv)
	log.Printf("Private key of node saved to key.private")

	return &Node{
		id:         &pub,
		listenPort: port,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

// load node with existing keys and settings
// func LoadNode() {
// 	priv := readKeyFromFile()
// }

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
