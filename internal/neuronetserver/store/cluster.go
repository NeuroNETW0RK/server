package store

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
)

type ICluster interface {
	Create(c context.Context, db *gorm.DB, cluster *model.Cluster) (err error)
	GetBy(c context.Context, db *gorm.DB, opts ...DBOptions) (cluster *model.ClusterBo, err error)
	GetListBy(c context.Context, db *gorm.DB, opts ...DBOptions) (clusters []model.ClusterBo, err error)
	GetCntBy(c context.Context, db *gorm.DB, opts ...DBOptions) (count int64, err error)
	DeleteBy(c context.Context, db *gorm.DB, opts ...DBOptions) (err error)
	UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...DBOptions) (err error)
	Updates(c context.Context, db *gorm.DB, cluster *model.Cluster, opts ...DBOptions) (err error)

	IClusterOption
}

type IClusterOption interface {
}
