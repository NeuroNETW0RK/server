package minio

import (
	"context"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// 分块上传 (partSize>=5M)
type BlockUpload struct {
	client          *Client
	minioClient     *minio.Client
	tmpBucketName   string
	uploadInfoSlice []*minio.UploadInfo
}

func NewBlockUpload(client *Client, tmpBucketName string) (*BlockUpload, error) {
	minioClient, err := newMinioClient(client)
	if err != nil {
		return nil, err
	}

	return &BlockUpload{
		client:        client,
		minioClient:   minioClient,
		tmpBucketName: tmpBucketName,
	}, nil
}

func (b *BlockUpload) Upload(ctx context.Context, r io.Reader, objectSize int64) error {
	objectName := strings.ReplaceAll(uuid.NewString(), "-", "")
	uploadInfo, err := b.minioClient.PutObject(ctx, b.tmpBucketName, objectName, r, objectSize, minio.PutObjectOptions{})
	if err == nil {
		b.uploadInfoSlice = append(b.uploadInfoSlice, &uploadInfo)
	}

	return err
}

func (b *BlockUpload) ComposeObject(ctx context.Context, objectname string, contentType string, metadata ...map[string]string) (*minio.UploadInfo, error) {
	srcs := make([]minio.CopySrcOptions, len(b.uploadInfoSlice))
	for index, uploadInfo := range b.uploadInfoSlice {
		srcs[index] = minio.CopySrcOptions{
			Bucket:    uploadInfo.Bucket,
			Object:    uploadInfo.Key,
			MatchETag: uploadInfo.ETag,
		}
	}

	userMetadata := map[string]string{}
	if contentType != "" {
		userMetadata["content-type"] = contentType
	}
	for _, meta := range metadata {
		for k, v := range meta {
			userMetadata[k] = v
		}
	}

	dst := minio.CopyDestOptions{
		Bucket:          b.client.Bucket,
		Object:          objectname,
		ReplaceMetadata: true,
		UserMetadata:    userMetadata,
	}

	uploadInfo, err := b.minioClient.ComposeObject(ctx, dst, srcs...)
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}

func (b *BlockUpload) RemoveTmpBlock(ctx context.Context) {
	for _, uploadInfo := range b.uploadInfoSlice {
		b.minioClient.RemoveObject(ctx, uploadInfo.Bucket, uploadInfo.Key, minio.RemoveObjectOptions{})
	}
}
