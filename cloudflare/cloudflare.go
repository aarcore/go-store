package cloudflare

import (
	"sync"

	"github.com/aarcore/go-store/store"
)

type Config struct {
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	AccountId string `mapstructure:"account_id" json:"account_id" yaml:"account_id"`
}

type Cloudflare struct {
	config *Config
}

var (
	cloudflare *Cloudflare
	once       *sync.Once
)

func Init(config Config) (store.Store, error) {
	return nil, nil
}
