package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
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
