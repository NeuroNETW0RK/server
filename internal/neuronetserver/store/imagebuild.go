package store

import (
	"context"
	"neuronet/internal/neuronetserver/model"

	"gorm.io/gorm"
)

type IImageBuild interface {
	Create(c context.Context, db *gorm.DB, image *model.ImageBuild) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (Tag *model.ImageBuild, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (imageBuilds []model.ImageBuild, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, imageBuild *model.ImageBuild, opts ...DBOptions) (err error)

	IImageBuildOption
}

type IImageBuildOption interface {
}
