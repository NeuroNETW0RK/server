package mysql

import (
	"context"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"

	"gorm.io/gorm"
)

var _ store.IImageTag = (*imageTag)(nil)

func newImageTag() *imageTag {
	return &imageTag{}
}

type imageTag struct {
}

func (s *imageTag) Create(c context.Context, db *gorm.DB, imageTag *model.ImageTag) (err error) {
	err = db.WithContext(c).Create(&imageTag).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (imageTag *model.ImageTag, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&imageTag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (imageTags []model.ImageTag, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&imageTags).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.ImageTag{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.ImageTag{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.ImageTag{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageTag) Updates(c context.Context, db *gorm.DB, cluster *model.ImageTag, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&cluster).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
