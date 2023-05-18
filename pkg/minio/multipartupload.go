package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// 多部分上传 (partSize>=5M)
type MultipartUpload struct {
	client     *Client
	minioCore  *minio.Core
	objectName string

	uploadID      string
	completeParts []minio.CompletePart
}

func NewMultipartUpload(client *Client, objectName string) (*MultipartUpload, error) {
	minioClient, err := newMinioClient(client)
	if err != nil {
		return nil, err
	}

	return &MultipartUpload{
		client:     client,
		minioCore:  &minio.Core{minioClient},
		objectName: objectName,
	}, nil
}

func (mpu *MultipartUpload) UploadPart(ctx context.Context, r io.Reader, partSize int64) error {
	if mpu.uploadID == "" {
		uploadID, err := mpu.minioCore.NewMultipartUpload(ctx, mpu.client.Bucket, mpu.objectName, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		mpu.uploadID = uploadID
	}
	part, err := mpu.minioCore.PutObjectPart(context.Background(), mpu.client.Bucket, mpu.objectName, mpu.uploadID, len(mpu.completeParts)+1, r, partSize, minio.PutObjectPartOptions{})
	if err != nil {
		return err
	}
	mpu.completeParts = append(mpu.completeParts, minio.CompletePart{PartNumber: part.PartNumber, ETag: part.ETag})

	return nil
}

func (mpu *MultipartUpload) Complate(ctx context.Context, contentType string) (minio.ObjectInfo, error) {
	if len(mpu.completeParts) == 0 {
		return minio.ObjectInfo{}, nil
	}

	if _, err := mpu.minioCore.CompleteMultipartUpload(context.Background(), mpu.client.Bucket, mpu.objectName, mpu.uploadID, mpu.completeParts, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return minio.ObjectInfo{}, err
	}

	return mpu.minioCore.StatObject(context.Background(), mpu.client.Bucket, mpu.objectName, minio.StatObjectOptions{})
}
