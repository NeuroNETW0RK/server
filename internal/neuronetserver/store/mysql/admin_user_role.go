package mysql

import (
	"NeuroNET/internal/neuronetserver/model"
	"NeuroNET/internal/neuronetserver/store"
	"NeuroNET/internal/pkg/code"
	"NeuroNET/pkg/errors"
	"context"
	"gorm.io/gorm"
)

var _ store.IUserRole = (*userRole)(nil)

func newUserRole() *userRole {
	return &userRole{}
}

type userRole struct {
}

func (u *userRole) CreateInBatch(c context.Context, db *gorm.DB, userRoles []model.UserRole) (err error) {
	err = db.WithContext(c).Create(&userRoles).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (u *userRole) Create(c context.Context, db *gorm.DB, userRole *model.UserRole) (err error) {
	err = db.WithContext(c).Create(&userRole).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (u *userRole) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (userRoles []model.UserRole, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&userRoles).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (u *userRole) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.UserRole{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (u *userRole) WithUserIDs(userIDs []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id IN ?", userIDs)
	}
}

func (u *userRole) WithRoleIDs(roleIDs []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role_id IN ?", roleIDs)
	}
}
