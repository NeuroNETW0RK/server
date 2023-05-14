package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IImage interface {
	Create(c context.Context, db *gorm.DB, image *model.Image) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (image *model.ImageBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (images []model.ImageBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, image *model.Image, opts ...DBOptions) (err error)

	IImageOption
}

type IImageOption interface {
}
