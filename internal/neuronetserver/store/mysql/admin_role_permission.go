package mysql

import (
	"NeuroNET/internal/neuronetserver/model"
	"NeuroNET/internal/neuronetserver/store"
	"NeuroNET/internal/pkg/code"
	"NeuroNET/pkg/errors"
	"context"
	"gorm.io/gorm"
)

var _ store.IRolePermission = (*rolePermission)(nil)

func newRolePermission() *rolePermission {
	return &rolePermission{}
}

type rolePermission struct {
}

func (r *rolePermission) Create(c context.Context, db *gorm.DB, permissionRoles []model.PermissionRole) (err error) {
	err = db.WithContext(c).Create(&permissionRoles).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (r *rolePermission) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (permissionRoles []model.PermissionRole, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&permissionRoles).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (r *rolePermission) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.PermissionRole{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (r *rolePermission) WithPermissionIDs(permissionIDs []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("permission_id IN ?", permissionIDs)
	}
}

func (r *rolePermission) WithRoleIDs(roleIDs []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role_id IN ?", roleIDs)
	}
}
