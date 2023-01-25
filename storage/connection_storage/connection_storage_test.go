package connection_storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionStorage(t *testing.T) {
	store := NewConnectionStorage("data/connections/")
	conns := []string{0: "1.1.1.1:3000", 1: "2.2.2.2:20000", 2: "3.3.3.3:55555"}

	saveerr := store.SaveConnections(conns)
	if saveerr != nil {
		panic(saveerr)
	}
	_, err := os.Stat("data/connections/connections.txt")

	assert.Nil(t, err)

	loaded := store.LoadConnections()

	assert.Equal(t, len(loaded), 3)
	assert.Equal(t, loaded[0], "1.1.1.1:3000")

	os.RemoveAll("data")
}
