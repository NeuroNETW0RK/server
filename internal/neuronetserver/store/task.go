package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type ITask interface {
	Create(c context.Context, db *gorm.DB, task *model.Task) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (task *model.TaskBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (tasks []model.TaskBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, task *model.Task, opts ...DBOptions) (err error)

	ITaskOption
}

type ITaskOption interface {
}
