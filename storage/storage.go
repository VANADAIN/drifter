package storage

import (
	"github.com/VANADAIN/drifter/storage/connection_storage"
	"github.com/VANADAIN/drifter/storage/keys_storage"
)

type StorageManager struct {
	ks *keys_storage.KeysStorage
	cs *connection_storage.ConnectionStorage
}

func NewStorageManager(keyPath string, connectionPath string) *StorageManager {
	keystore := keys_storage.NewKeysStorage(keyPath)
	connectionstore := connection_storage.NewConnectionStorage(connectionPath)

	sm := &StorageManager{
		ks: keystore,
		cs: connectionstore,
	}

	return sm
}
