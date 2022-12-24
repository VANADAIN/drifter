package node

import (
	"log"
	"net"
	"os"

	"github.com/VANADAIN/drifter/dcrypto"
)

type Node struct {
	id         *dcrypto.PublicKey
	listenPort string
	lsn        net.Listener
	msgch      chan Message
	quitch     chan struct{}
}

func NewNode(port string) *Node {
	priv := dcrypto.GeneratePrivateKey()
	pub := priv.Public()

	status, err := exists("./keys")
	if err != nil {
		panic("Error reading keys folder")
	}

	if !status {
		err := os.Mkdir("./keys", 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		panic("Wrong entrypoint for node! Switch to Load mode.")
	}

	saveKeyToFile(&priv.Key)
	log.Printf("Private key of node saved to key.private")

	// TODO: read port from settings file
	return &Node{
		id:         pub,
		listenPort: port,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

// load node with existing keys and settings
func LoadNode() *Node {
	rawpriv := readKeyFromFile()
	priv := dcrypto.PrivateKey{
		Key: rawpriv,
	}

	pub := priv.Public()

	return &Node{
		id:         pub,
		listenPort: ":3000",
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}
