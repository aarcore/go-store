package store

import (
	"fmt"
	"io"
)

type Store interface {
	Put(key string, r io.Reader, dataLength int64) error
	PutFile(key string, localFile string) error
    Get(key string) (io.ReadCloser, error)
    Rename(srcKey string, destKey string) error
    Copy(srcKey string, destKey string) error
    Exists(key string) (bool, error)
    Size(key string) (int64, error)
    Delete(key string) error
    Url(key string) string
}

var stores = make(map[StoreName]Store)

func SetStore(name StoreName, store Store) {
	if store == nil {
		panic("store: Register disk is nil")
	}

	stores[name] = store
}

func GetStore(name StoreName) (Store, error) {
	store, err := stores[name]
	if !err {
		return nil, fmt.Errorf("Unknown store %q", name)
	}

	return store, nil
}