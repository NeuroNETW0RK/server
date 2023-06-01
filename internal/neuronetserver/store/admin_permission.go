package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IPermission interface {
	Create(c context.Context, db *gorm.DB, permission *model.Permission) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (permission *model.PermissionDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (permissions []model.PermissionDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, permission *model.Permission, opts ...DBOptions) (err error)

	IPermissionOption
}

type IPermissionOption interface {
	WithResource(resource string) DBOptions
}
