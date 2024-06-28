package ali

import (
	"io"
	"strconv"
	"sync"

	"github.com/aarcore/go-store/store"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Config struct {
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	IsSSL           bool   `mapstructure:"is_ssl" json:"is_ssl" yaml:"is_ssl"`
	IsPrivate       bool   `mapstructure:"is_private" json:"is_private" yaml:"is_private"`
}

type ali struct {
	config *Config
	client *alioss.Client
	bucket *alioss.Bucket
}

var (
	o       *ali
	once    *sync.Once
	initErr error
)

func Init(config Config) (store.Store, error) {
	once = &sync.Once{}
	once.Do(func() {
		o = &ali{}
		o.config = &config
		o.client, initErr = alioss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
		if initErr != nil {
			return
		}

		o.bucket, initErr = o.client.Bucket(config.Bucket)
		if initErr != nil {
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return o, nil
}


func (a *ali) Copy(srcKey string, destKey string) error {
	srcKey = store.NormalizeKey(srcKey)
	destKey = store.NormalizeKey(destKey)

	_, err := o.bucket.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}

	return nil
}

func (a *ali) Delete(key string) error {
	key = store.NormalizeKey(key)

    err := o.bucket.DeleteObject(key)
    if err != nil {
        return err
    }

    return nil
}

func (a *ali) Exists(key string) (bool, error) {
	key = store.NormalizeKey(key)

	return o.bucket.IsObjectExist(key)
}

func (a *ali) Get(key string) (io.ReadCloser, error) {
	key = store.NormalizeKey(key)

	body, err := o.bucket.GetObject(key)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (a *ali) Put(key string, r io.Reader, dataLength int64) error {
	key = store.NormalizeKey(key)

	err := o.bucket.PutObject(key, r)
	if err != nil {
		return err
	}

	return nil
}

func (a *ali) PutFile(key string, localFile string) error {
	key = store.NormalizeKey(key)

	err := o.bucket.PutObjectFromFile(key, localFile)
	if err != nil {
		return err
	}

	return nil
}

func (a *ali) Rename(srcKey string, destKey string) error {
	srcKey = store.NormalizeKey(srcKey)
	destKey = store.NormalizeKey(destKey)

	_, err := o.bucket.CopyObject(srcKey, destKey)
	if err != nil {
		return err
	}

	err = o.Delete(srcKey)
	if err != nil {
		return err
	}

	return nil
}

func (a *ali) Size(key string) (int64, error) {
	key = store.NormalizeKey(key)

	props, err := o.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return 0, err
	}

	size, err := strconv.ParseInt(props.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (a *ali) Url(key string) string {
	var prefix string
	key = store.NormalizeKey(key)

	if o.config.IsSSL {
		prefix = "https://"
	} else {
		prefix = "http://"
	}

	if o.config.IsPrivate {
		url, err := o.bucket.SignURL(key, alioss.HTTPGet, 3600)
		if err == nil {
			return url
		}
	}

	return prefix + o.config.Bucket + "." + o.config.Endpoint + "/" + key
}

