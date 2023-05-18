package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (store *Client) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return nil, err
	}

	bucketInfoList, err := minioClient.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return bucketInfoList, nil
}
