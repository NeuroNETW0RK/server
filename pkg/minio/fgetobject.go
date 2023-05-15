package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (store *Client) DownloadAOSObject(ctx context.Context, remotePath, localPath string) (err error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return err
	}

	err = minioClient.FGetObject(ctx, store.Bucket, remotePath, localPath, minio.GetObjectOptions{})
	return err
}
