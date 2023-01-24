package connection_storage

import (
	"bufio"
	"fmt"
	"os"
)

type ConnectionStorage struct {
	FilePath string
}

func NewConnectionStorage(path string) *ConnectionStorage {
	store := &ConnectionStorage{
		FilePath: path,
	}

	return store
}

func (cs *ConnectionStorage) SaveConnections(conns []string) error {
	cs.DeleteFile()
	f, err := os.Create(cs.FilePath)
	if err != nil {
		panic("Error trying to create known connections file")
	}

	defer f.Close()
	if err != nil {
		panic("Error trying to save known connections file")
	}

	for _, value := range conns {
		fmt.Fprintln(f, value)
	}

	return nil
}

func (cs *ConnectionStorage) LoadConnections() []string {
	res := make([]string, 0)
	file, err := os.Open(cs.FilePath)
	if err != nil {
		panic("Error reading connections file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		res = append(res, line)
	}

	return res
}

func (cs *ConnectionStorage) DeleteFile() {
	err := os.Remove(cs.FilePath)
	if err != nil {
		panic("Error trying to delete known connections file")
	}
}
