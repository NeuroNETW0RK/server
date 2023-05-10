package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type ITaskImage interface {
	Create(c context.Context, db *gorm.DB, taskImage *model.TaskImage) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (taskImage *model.TaskImageBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (taskImages []model.TaskImageBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, taskImage *model.TaskImage, opts ...DBOptions) (err error)

	ITaskImageOption
}

type ITaskImageOption interface {
}
