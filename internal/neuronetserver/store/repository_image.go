package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IRepositoryImage interface {
	Create(c context.Context, db *gorm.DB, repositoryImage *model.RepositoryImage) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repositoryImage *model.RepositoryImageDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repositoryImages []model.RepositoryImageDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, repositoryImage *model.RepositoryImage, opts ...DBOptions) (err error)

	IRepositoryImageOption
}

type IRepositoryImageOption interface {
}
