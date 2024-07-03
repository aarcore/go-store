package r2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/aarcore/go-store/store"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Config struct {
	AccountId string `mapstructure:"account_id" json:"account_id" yaml:"account_id"`
	AccessKey string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}

type R2 struct {
	r2config *R2Config
	s3client *s3.Client
	bucket   string
}

var (
	r2  *R2
	once *sync.Once
)

func Init(r2config R2Config) (store.Store, error) {
	once = &sync.Once{}
	once.Do(func() {
		r2 = &R2{}
		r2.r2config = &r2config

		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2.r2config.AccountId),
			}, nil
		})

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolverWithOptions(resolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2.r2config.AccessKey, r2.r2config.SecretKey, "")),
			config.WithRegion("auto"),
		)

		r2.s3client = s3.NewFromConfig(cfg)
		// output, err := r2.s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		// 	Bucket: &r2.r2config.Bucket,
		// })

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(output)

		store.SetStore(store.R2, r2)
	})

	return r2, nil
}

// Copy implements store.Store.
func (r *R2) Copy(srcKey string, destKey string) error {
	panic("unimplemented")
}

// Delete implements store.Store.
func (r *R2) Delete(key string) error {
	output, err := r2.s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &r2.r2config.Bucket,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, object := range output.Contents {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

	return nil
}

func (r *R2) Exists(key string) (bool, error) {
	panic("unimplemented")
}

func (r2 *R2) Get(key string) (io.ReadCloser, error) {
	output, err := r2.s3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Key: aws.String(key),
		Bucket: aws.String(r2.r2config.Bucket),
	})

	if err != nil {
		return nil, err
	}

	defer output.Body.Close()

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return string(body), nil
}

// Put implements store.Store.
func (r2 *R2) Put(key string, r io.Reader, dataLength int64) error {
	_, err := r2.s3client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key: aws.String(r2.r2config.AccessKey),
		Bucket: aws.String(r2.r2config.Bucket),
		Body: strings.NewReader("Hello, r2"),
	})

	if err != nil {
		return err
	}

	return nil
}

// PutFile implements store.Store.
func (r *R2) PutFile(key string, localFile string) error {
	panic("unimplemented")
}

// Rename implements store.Store.
func (r *R2) Rename(srcKey string, destKey string) error {
	panic("unimplemented")
}

// Size implements store.Store.
func (r *R2) Size(key string) (int64, error) {
	panic("unimplemented")
}

// Url implements store.Store.
func (r *R2) Url(key string) string {
	panic("unimplemented")
}