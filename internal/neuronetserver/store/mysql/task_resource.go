package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.ITaskResource = (*taskResource)(nil)

func newTaskResource() *taskResource {
	return &taskResource{}
}

type taskResource struct {
}

func (s *taskResource) Create(c context.Context, db *gorm.DB, taskResource *model.TaskResource) (err error) {
	err = db.WithContext(c).Create(&taskResource).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (taskResource *model.TaskResourceDo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&taskResource).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (taskResources []model.TaskResourceDo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&taskResources).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.TaskResource{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.TaskResource{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.TaskResource{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *taskResource) Updates(c context.Context, db *gorm.DB, taskResource *model.TaskResource, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&taskResource).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
