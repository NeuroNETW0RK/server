package store

import (
	"context"
	"neuronet/internal/neuronetserver/model"

	"gorm.io/gorm"
)

type IImage interface {
	Create(c context.Context, db *gorm.DB, image *model.ImageInfo) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (image *model.ImageInfo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (images []model.ImageInfo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, image *model.ImageInfo, opts ...DBOptions) (err error)

	IImageOption
}

type IImageOption interface {
}
