package store

import (
	"NeuroNET/internal/neuronetserver/model"
	"context"
	"gorm.io/gorm"
)

type IUser interface {
	Create(c context.Context, db *gorm.DB, user *model.User) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (user *model.UserBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (user []model.UserBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, user *model.User, opts ...DBOptions) (err error)

	IUserOption
}

type IUserOption interface {
	WithAccount(account string) DBOptions
}
