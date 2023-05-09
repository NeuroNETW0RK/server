package mysql

import (
	"NeuroNET/internal/neuronetserver/model"
	"NeuroNET/internal/neuronetserver/store"
	"NeuroNET/internal/pkg/code"
	"NeuroNET/pkg/errors"
	"context"
	"gorm.io/gorm"
)

var _ store.IPermission = (*permission)(nil)

func newPermission() *permission {
	return &permission{}
}

type permission struct {
}

func (s *permission) Updates(c context.Context, db *gorm.DB, permission *model.Permission, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&permission).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) Create(c context.Context, db *gorm.DB, permission *model.Permission) (err error) {
	err = db.WithContext(c).Create(&permission).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (permission *model.PermissionBo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (permissions []model.PermissionBo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&permissions).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.Permission{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.Permission{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.Permission{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *permission) WithResource(resource string) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("resource = ?", resource)
	}
}
