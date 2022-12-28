package node

import (
	"github.com/VANADAIN/drifter/network"
)

type Node struct {
	nmng *network.NetworkManager
}

// func NewNode(port string) *Node {
// 	priv := dcrypto.GeneratePrivateKey()
// 	pub := priv.Public()

// 	settings := settings.Read()
// 	keys_path := settings.Keys_path

// 	checkKeysSaved(keys_path)
// 	fullpath := keys_path + "/key.private"

// 	saveKeyToFile(priv, fullpath)
// 	log.Printf("Private key of node saved to key.private")

// 	// TODO: create network manager

// 	// return &Node{
// 	// 	id: pub,
// 	// }
// }

// load node with existing keys and settings
// func LoadNode() *Node {
// 	priv := readKeyFromFile("./keys_saved/key.private")
// 	pub := priv.Public()

// 	return &Node{
// 		id: pub,
// 	}

// }
