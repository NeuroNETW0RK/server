package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type ITaskResource interface {
	Create(c context.Context, db *gorm.DB, taskResource *model.TaskResource) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (taskResource *model.TaskResourceDo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (taskResources []model.TaskResourceDo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, taskResource *model.TaskResource, opts ...DBOptions) (err error)

	ITaskResourceOption
}

type ITaskResourceOption interface {
}
