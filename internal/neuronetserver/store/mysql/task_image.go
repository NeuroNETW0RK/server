package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.ITaskImage = (*taskImage)(nil)

func newTaskImage() *taskImage {
	return &taskImage{}
}

type taskImage struct {
}

func (s *taskImage) Create(c context.Context, db *gorm.DB, taskImage *model.TaskImage) (err error) {
	err = db.WithContext(c).Create(&taskImage).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (taskImage *model.TaskImageBo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&taskImage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (taskImages []model.TaskImageBo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&taskImages).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.TaskImage{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.TaskImage{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.TaskImage{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskImage) Updates(c context.Context, db *gorm.DB, taskImage *model.TaskImage, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&taskImage).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
