package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.IRepositoryImage = (*repositoryImage)(nil)

func newRepositoryImage() *repositoryImage {
	return &repositoryImage{}
}

type repositoryImage struct {
}

func (s *repositoryImage) Create(c context.Context, db *gorm.DB, repositoryImage *model.RepositoryImage) (err error) {
	err = db.WithContext(c).Create(&repositoryImage).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (repositoryImage *model.RepositoryImageDo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&repositoryImage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (repositoryImages []model.RepositoryImageDo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&repositoryImages).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.RepositoryImage{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.RepositoryImage{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.RepositoryImage{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImage) Updates(c context.Context, db *gorm.DB, repositoryImage *model.RepositoryImage, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&repositoryImage).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
