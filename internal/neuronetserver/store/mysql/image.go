package mysql

import (
	"context"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"

	"gorm.io/gorm"
)

var _ store.IImage = (*image)(nil)

func newImage() *image {
	return &image{}
}

type image struct {
}

func (s *image) Create(c context.Context, db *gorm.DB, image *model.ImageInfo) (err error) {
	err = db.WithContext(c).Create(&image).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (image *model.ImageInfo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (images []model.ImageInfo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&images).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.ImageInfo{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.ImageInfo{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.ImageInfo{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *image) Updates(c context.Context, db *gorm.DB, imageInfo *model.ImageInfo, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&imageInfo).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
