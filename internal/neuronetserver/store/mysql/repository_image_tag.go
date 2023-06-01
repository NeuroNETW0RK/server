package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.IRepositoryImageTag = (*repositoryImageTag)(nil)

func newRepositoryImageTag() *repositoryImageTag {
	return &repositoryImageTag{}
}

type repositoryImageTag struct {
}

func (s *repositoryImageTag) Create(c context.Context, db *gorm.DB, repositoryImageTag *model.RepositoryImageTag) (err error) {
	err = db.WithContext(c).Create(&repositoryImageTag).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (repositoryImageTag *model.RepositoryImageTagDo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&repositoryImageTag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (repositoryImageTags []model.RepositoryImageTagDo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&repositoryImageTags).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.RepositoryImageTag{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.RepositoryImageTag{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.RepositoryImageTag{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *repositoryImageTag) Updates(c context.Context, db *gorm.DB, repositoryImageTag *model.RepositoryImageTag, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&repositoryImageTag).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
