package node

import (
	"log"
	"net"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/settings"
)

type Node struct {
	id          *dcrypto.PublicKey
	ListenPort  string
	Lsn         net.Listener
	ActiveConns []net.Conn
	ConnList    map[string]string // strings of pubk -> net.Addr.String()
	Aliases     map[string]string // name (alias) -> pubk
	Msgch       chan Message
	Quitch      chan struct{}
}

func NewNode(port string) *Node {
	priv := dcrypto.GeneratePrivateKey()
	pub := priv.Public()

	settings := settings.Read()
	keys_path := settings.Keys_path

	checkKeysSaved(keys_path)
	fullpath := keys_path + "/key.private"

	saveKeyToFile(priv, fullpath)
	log.Printf("Private key of node saved to key.private")

	return &Node{
		id:         pub,
		ListenPort: port,
		Quitch:     make(chan struct{}),
		Msgch:      make(chan Message, 10),
	}
}

// load node with existing keys and settings
func LoadNode() *Node {
	priv := readKeyFromFile("./keys_saved/key.private")
	pub := priv.Public()

	return &Node{
		id:         pub,
		ListenPort: ":3000",
		Quitch:     make(chan struct{}),
		Msgch:      make(chan Message, 10),
	}
}
