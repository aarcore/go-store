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

func (qiniu *Qiniu) Put(key string, r io.Reader, dataLength int64) error {
	key = store.NormalizeKey(key)
	upToken := qiniu.putPolicy.UploadToken(qiniu.mac)
	ret := qiniuStore.PutRet{}
	err := qiniu.formUploader.Put(context.Background(), &ret, upToken, key, r, dataLength, nil)

	if err != nil {
		return err
	}

	return nil
}

func (qiniu *Qiniu) PutFile(key string, localFile string) error {
	key = store.NormalizeKey(key)
	upToken := qiniu.putPolicy.UploadToken(qiniu.mac)
	ret := qiniuStore.PutRet{}
	err := qiniu.formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		return err
	}

	return nil
}

func (qiniu *Qiniu) Copy(souKey string, destKey string) error {
	souKey = store.NormalizeKey(souKey)
	destKey = store.NormalizeKey(destKey)

	err := qiniu.bucketManager.Copy(qiniu.config.Bucket, souKey, qiniu.config.Bucket, destKey, true)
	if err != nil {
		return err
	}

	return nil
}

func (qiniu *Qiniu) Delete(key string) error {
	key = store.NormalizeKey(key)
	err := qiniu.bucketManager.Delete(qiniu.config.Bucket, key)
	if err != nil {
		return err
	}

	return nil
}

func (qiniu *Qiniu) Exists(key string) (bool, error) {
	return true, nil
}

func (qiniu *Qiniu) Get(key string) (io.ReadCloser, error) {
	key = store.NormalizeKey(key)
	res, err := http.Get(qiniu.Url(key))
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (qiniu *Qiniu) Rename(sourceKey string, destKey string) error {
	sourceKey = store.NormalizeKey(sourceKey)
	destKey = store.NormalizeKey(destKey)

	err := qiniu.bucketManager.Move(qiniu.config.Bucket, sourceKey, qiniu.config.Bucket, destKey, true)
	if err != nil {
		return err
	}

	return nil
}

func (qiniu *Qiniu) Size(key string) (int64, error) {
	key = store.NormalizeKey(key)
	info, err := qiniu.bucketManager.Stat(qiniu.config.Bucket, key)
	if err != nil {
		return 0, err
	}
	return info.Fsize, nil
}

func (qiniu *Qiniu) Url(key string) string{
	var prefix string
	key = store.NormalizeKey(key)

	if qiniu.config.IsSSL {
		prefix = "https://"
	} else {
		prefix = "http://"
	}

	if qiniu.config.IsPrivate {
		deadline := time.Now().Add(time.Second * 3600).Unix()
		return prefix + qiniuStore.MakePrivateURL(qiniu.mac, qiniu.config.Domain, key, deadline)
	}

	return prefix + qiniuStore.MakePublicURL(qiniu.config.Domain, key)
}
