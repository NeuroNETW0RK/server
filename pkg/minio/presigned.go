package minio

import (
	"context"
	"net/url"
	"time"
)

const (
	PreSignedExpiresOneHour = time.Hour
	PreSignedExpiresOneDay  = time.Hour * 24
	PreSignedExpiresMin     = time.Second
	PreSignedExpiresMax     = time.Hour * 24 * 7 // store 对象存储获取URL最长时间为7天
)

func (store *Client) GetAOSObjURL(ctx context.Context, filePath string, expires time.Duration) (u *url.URL, err error) {
	if expires > PreSignedExpiresMax {
		expires = PreSignedExpiresMax
	}
	if expires < PreSignedExpiresMin {
		expires = PreSignedExpiresMin
	}
	minioClient, err := newMinioClient(store)
	if err != nil {
		return nil, err
	}

	val := make(url.Values)
	u, err = minioClient.PresignedGetObject(ctx, store.Bucket, filePath, expires, val)

	return
}
