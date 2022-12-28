package node

import (
	"encoding/json"
	"log"
	"os"

	"github.com/VANADAIN/drifter/dcrypto"
)

func checkKeysSaved(path string) {
	status, err := exists(path)
	if err != nil {
		panic("Error reading keys folder")
	}

	if !status {
		err := os.Mkdir("./keys_saved", 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		panic("Wrong entrypoint for node! Switch to Load mode.")
	}
}

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
func saveKeyToFile(key *dcrypto.PrivateKey, path string) error {
	err := os.WriteFile(path, key.Key, 0666)
	if err != nil {
		log.Fatal(err)
		panic("Error saving keys to file")
	}

	return nil
}

// for node loading
func readKeyFromFile(path string) *dcrypto.PrivateKey {
	key_raw, err := os.ReadFile(path)
	if err != nil {
		panic("Error reading keys file")
	}

	key := &dcrypto.PrivateKey{
		Key: key_raw,
	}

	return key
}

func saveConnectionList(n *Node) {
	jsonBytes, err := json.MarshalIndent(n.nmng.ConnList, "", "    ")
	if err != nil {
		log.Fatal("Error converting connection list to json")
	}

	werr := os.WriteFile("./nodeinfo/coonectionlist.json", jsonBytes, 0666)
	if err != nil {
		log.Fatal(werr)
		panic("Error saving keys to file")
	}
}

func loadConnectionList(n *Node) {
	b, err := os.ReadFile("./nodeinfo/connectionlist.json")
	if err != nil {
		panic("Error downloading connection list")
	}

	jerr := json.Unmarshal(b, &n.nmng.ConnList)
	if jerr != nil {
		panic("Error unmarshaling connection list")
	}
}
