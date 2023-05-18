package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (store *Client) GetAOSObject(ctx context.Context, filePath string) (obj *minio.Object, err error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return nil, err
	}

	obj, err = minioClient.GetObject(ctx, store.Bucket, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return
}
