package store

import (
	"NeuroNET/internal/neuronetserver/model"
	"context"
	"gorm.io/gorm"
)

type IRole interface {
	Create(c context.Context, db *gorm.DB, role *model.Role) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (role *model.RoleBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (roles []model.RoleBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, role *model.Role, opts ...DBOptions) (err error)

	IRoleOption
}

type IRoleOption interface {
}
