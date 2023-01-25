package connection_storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	if _, err := os.Stat(cs.FilePath); err == nil {
		// mixing is done in memory, so we can delete file if file exists
		cs.DeleteFile()
	}

	cs.CreatePath()

	f, err := os.Create(cs.FilePath + "connections.txt")
	if err != nil {
		fmt.Println(err)
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
	file, err := os.Open(cs.FilePath + "connections.txt")
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
	err := os.Remove(cs.FilePath + "connections.txt")
	if err != nil {
		panic("Error trying to delete known connections file")
	}
}

func (cs *ConnectionStorage) CreatePath() {
	// create path
	newpath := filepath.Join(".", cs.FilePath)
	e := os.MkdirAll(newpath, os.ModePerm)
	if e != nil {
		panic(e)
	}
}
