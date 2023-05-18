package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
)

func (store *Client) GetObjectTags(ctx context.Context, objectName string) (map[string]string, error) {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return nil, err
	}

	tags, err := minioClient.GetObjectTagging(ctx, store.Bucket, objectName, minio.GetObjectTaggingOptions{})
	if err != nil {
		return nil, err
	}

	if tags.TagSet == nil {
		return nil, nil
	}

	return tags.ToMap(), nil
}

func (store *Client) PutObjectTags(ctx context.Context, objectName string, tagMap map[string]string) error {
	minioClient, err := newMinioClient(store)
	if err != nil {
		return err
	}

	tagging, err := tags.NewTags(tagMap, true)
	if err != nil {
		return err
	}

	return minioClient.PutObjectTagging(ctx, store.Bucket, objectName, tagging, minio.PutObjectTaggingOptions{})
}
