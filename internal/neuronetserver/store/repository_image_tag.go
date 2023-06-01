package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IRepositoryImageTag interface {
	Create(c context.Context, db *gorm.DB, repositoryImageTag *model.RepositoryImageTag) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repositoryImageTag *model.RepositoryImageTagDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repositories []model.RepositoryImageTagDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, repositoryImageTag *model.RepositoryImageTag, opts ...DBOptions) (err error)

	IRepositoryImageTagOption
}

type IRepositoryImageTagOption interface {
}
