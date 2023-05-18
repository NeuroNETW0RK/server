package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (store *Client) StatObject(ctx context.Context, filePath string) (objInfo minio.ObjectInfo, err error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return minio.ObjectInfo{}, err
	}

	objInfo, err = minioClient.StatObject(ctx, store.Bucket, filePath, minio.StatObjectOptions{})
	return
}
