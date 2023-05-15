package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

func (store *Client) Upload(ctx context.Context, objectName, filePath, contentType string, metadata ...map[string]string) (minio.UploadInfo, error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	options := makePutObjectOptions(contentType, metadata...)

	// 使用FPutObject上传一个zip文件。
	return minioClient.FPutObject(
		ctx, store.Bucket, objectName, filePath, *options)
}

func (store *Client) UploadStream(ctx context.Context, objectName string, r io.Reader, objectSize int64, contentType string, metadata ...map[string]string) (minio.UploadInfo, error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	options := makePutObjectOptions(contentType, metadata...)
	options.DisableMultipart = objectSize > 0

	// 使用PutObject上传一个reader流。
	return minioClient.PutObject(ctx, store.Bucket, objectName, r, objectSize, *options)
}

func makePutObjectOptions(contentType string, metadata ...map[string]string) *minio.PutObjectOptions {
	options := minio.PutObjectOptions{
		ContentType:      contentType,
		UserMetadata:     make(map[string]string),
		DisableMultipart: true,
	}

	for _, meta := range metadata {
		for k, v := range meta {
			options.UserMetadata[k] = v
		}
	}

	return &options
}
