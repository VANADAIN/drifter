package server

import (
	"errors"
	"os"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/nat"
	"github.com/VANADAIN/drifter/storage"
)

const DEF_NODE_INFO_PATH = "node_info/node_info.json"

type Server struct {
	Name string
	IP   string
	ID   dcrypto.PublicKey
	CH   *ConnectionHandler
	SM   *storage.StorageManager
}

func CheckInitialized() bool {
	// check if server has some info saved
	if _, err := os.Stat(DEF_NODE_INFO_PATH); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func CreateInfo(name string) {
	// create keys and other, then save
}

func SaveInfo() {
	// save info to json
}

func LoadInfo() {
	// load info from json
}

func NewServer(name string) *Server {
	status := CheckInitialized()
	if !status {
		CreateInfo(name)
	}

	// now we're sure info exists or created
	// get info from load inf0
	LoadInfo()

	server := &Server{
		Name: name,
		IP:   nat.GetLocalIP().String(),
	}

	go server.CH.RunConnectionLoop()

	return server
}
