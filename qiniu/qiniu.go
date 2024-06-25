package qiniu

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/aarcore/go-store/store"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniuStore "github.com/qiniu/go-sdk/v7/storage"
)

type Config struct {
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain string `mapstructure:"domain" json:"domain" yaml:"domain"`
	IsSSL bool `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate bool `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}

type Qiniu struct {
	config *Config
	mac *qbox.Mac
	putPolicy *qiniuStore.PutPolicy
	formUploader *qiniuStore.FormUploader
	bucketManager *qiniuStore.BucketManager
}

var (
	qiniu *Qiniu
	once *sync.Once
)

func Init(config Config) (store.Store, error) {
	once = &sync.Once{}
	once.Do(func() {
		qiniu = &Qiniu{}
		qiniu.config = &config

		qiniu.putPolicy = &qiniuStore.PutPolicy{
			Scope: config.Bucket,
		}
		qiniu.mac = qbox.NewMac(config.AccessKey, config.SecretKey)

		cfg := qiniuStore.Config{
			UseHTTPS: config.IsSSL,
			UseCdnDomains: false,
		}

		qiniu.formUploader = qiniuStore.NewFormUploader(&cfg)
		qiniu.bucketManager = qiniuStore.NewBucketManager(qiniu.mac, &cfg)

		store.SetStore(store.Qiniu, qiniu)
	})

	return qiniu, nil
}
