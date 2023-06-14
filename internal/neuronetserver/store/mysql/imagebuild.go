package mysql

import (
	"context"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"

	"gorm.io/gorm"
)

var _ store.IImageBuild = (*imageBuild)(nil)

func newImageBuild() *imageBuild {
	return &imageBuild{}
}

type imageBuild struct {
}

func (s *imageBuild) Create(c context.Context, db *gorm.DB, image *model.ImageBuild) (err error) {
	err = db.WithContext(c).Create(&image).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (image *model.ImageBuild, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (images []model.ImageBuild, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&images).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.ImageBuild{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.ImageBuild{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.ImageBuild{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *imageBuild) Updates(c context.Context, db *gorm.DB, imageBuild *model.ImageBuild, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&imageBuild).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
