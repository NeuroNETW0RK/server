package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (store *Client) RenameObject(ctx context.Context, objectName, newBucket, newObjectName string) (minio.UploadInfo, error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	src := minio.CopySrcOptions{
		Bucket: store.Bucket,
		Object: objectName,
	}

	dest := minio.CopyDestOptions{
		Bucket: newBucket,
		Object: newObjectName,
	}

	uploadInfo, err := minioClient.CopyObject(ctx, dest, src)
	if err == nil {
		minioClient.RemoveObject(ctx, src.Bucket, src.Object, minio.RemoveObjectOptions{})
	}

	return uploadInfo, err
}
