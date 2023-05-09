package store

import (
	"NeuroNET/internal/neuronetserver/model"
	"context"
	"gorm.io/gorm"
)

type IRolePermission interface {
	Create(c context.Context, db *gorm.DB, permissionRoles []model.PermissionRole) (err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (permissionRoles []model.PermissionRole, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)

	IRolePermissionOption
}

type IRolePermissionOption interface {
	WithPermissionIDs(permissionIDs []int64) DBOptions
	WithRoleIDs(roleIDs []int64) DBOptions
}
