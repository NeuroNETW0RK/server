package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type IUserRole interface {
	Create(c context.Context, db *gorm.DB, userRole *model.UserRole) (err error)
	CreateInBatch(c context.Context, db *gorm.DB, userRoles []model.UserRole) (err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (userRoles []model.UserRole, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)

	IUserRoleOption
}

type IUserRoleOption interface {
	WithUserIDs(userIDs []int64) DBOptions
	WithRoleIDs(roleIDs []int64) DBOptions
}
