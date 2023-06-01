package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IRepository interface {
	Create(c context.Context, db *gorm.DB, repository *model.Repository) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repository *model.RepositoryDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (repositories []model.RepositoryDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, repository *model.Repository, opts ...DBOptions) (err error)

	IRepositoryOption
}

type IRepositoryOption interface {
}
