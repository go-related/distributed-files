package server

import "fmt"

type InMemoryFileStorage struct {
}

func NewInMemoryFileStorage() *InMemoryFileStorage {
	return &InMemoryFileStorage{}
}

func (fs *InMemoryFileStorage) Save(sc []byte) error {
	fmt.Println(fmt.Sprintf("storing file: %s", string(sc)))
	return nil
}
