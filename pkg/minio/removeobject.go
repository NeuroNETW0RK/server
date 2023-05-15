package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
)

func (store *Client) RemoveObject(ctx context.Context, objectName string) error {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return err
	}

	newCTX, _ := context.WithTimeout(ctx, 3*time.Second)
	return minioClient.RemoveObject(newCTX, store.Bucket, objectName, minio.RemoveObjectOptions{ForceDelete: true, GovernanceBypass: true})
}
