package store

import (
	"context"
	"neuronet/internal/neuronetserver/model"

	"gorm.io/gorm"
)

type IImageTag interface {
	Create(c context.Context, db *gorm.DB, image *model.ImageTag) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (Tag *model.ImageTag, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (imageTags []model.ImageTag, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, imageTag *model.ImageTag, opts ...DBOptions) (err error)

	IImageTagOption
}

type IImageTagOption interface {
}
