package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Host      string // "ip + port"
	AccessKey string // 访问用户
	Secret    string // 访问密钥
	Bucket    string // 桶名称
}

func newMinioClient(store *Client) (*minio.Client, error) {
	minioClient, err := minio.New(store.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(store.AccessKey, store.Secret, ""),
		Secure: false,
	})
	if err != nil {
		return &minio.Client{}, err
	}
	return minioClient, err
}
