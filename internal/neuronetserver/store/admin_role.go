package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IRole interface {
	Create(c context.Context, db *gorm.DB, role *model.Role) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (role *model.RoleDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (roles []model.RoleDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, role *model.Role, opts ...DBOptions) (err error)

	IRoleOption
}

type IRoleOption interface {
}
