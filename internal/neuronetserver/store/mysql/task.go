package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.ITask = (*task)(nil)

func newTask() *task {
	return &task{}
}

type task struct {
}

func (s *task) Create(c context.Context, db *gorm.DB, task *model.Task) (err error) {
	err = db.WithContext(c).Create(&task).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (task *model.TaskDo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (tasks []model.TaskDo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&tasks).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.Task{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.Task{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.Task{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *task) Updates(c context.Context, db *gorm.DB, task *model.Task, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&task).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
