package local

import (
	"io"
	"sync"

	"github.com/aarcore/go-store/store"
)

type Config struct {
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}

type local struct {
	config *Config
}

var (
	l    *local
	once *sync.Once
)

func Init(config Config) (store.Store, error) {
	once = &sync.Once{}
	once.Do(func() {
		l = &local{
			config: &config,
		}

		store.SetStore(store.Local, l)
	})

	return l, nil
}

// Copy implements store.Store.
func (l *local) Copy(srcKey string, destKey string) error {
	panic("unimplemented")
}

// Delete implements store.Store.
func (l *local) Delete(key string) error {
	panic("unimplemented")
}

// Exists implements store.Store.
func (l *local) Exists(key string) (bool, error) {
	panic("unimplemented")
}

// Get implements store.Store.
func (l *local) Get(key string) (io.ReadCloser, error) {
	panic("unimplemented")
}

// Put implements store.Store.
func (l *local) Put(key string, r io.Reader, dataLength int64) error {
	panic("unimplemented")
}

// PutFile implements store.Store.
func (l *local) PutFile(key string, localFile string) error {
	panic("unimplemented")
}

// Rename implements store.Store.
func (l *local) Rename(srcKey string, destKey string) error {
	panic("unimplemented")
}

// Size implements store.Store.
func (l *local) Size(key string) (int64, error) {
	panic("unimplemented")
}

// Url implements store.Store.
func (l *local) Url(key string) string {
	panic("unimplemented")
}


