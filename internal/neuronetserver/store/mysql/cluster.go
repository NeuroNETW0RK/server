package mysql

import (
	"context"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ store.ICluster = (*cluster)(nil)

func newCluster() *cluster {
	return &cluster{}
}

type cluster struct {
}

func (s *cluster) Create(c context.Context, db *gorm.DB, cluster *model.Cluster) (err error) {
	err = db.WithContext(c).Create(&cluster).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) GetBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (cluster *model.ClusterDo, err error) {
	err = db.WithContext(c).Scopes(opts...).First(&cluster).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) GetListBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (clusters []model.ClusterDo, err error) {
	err = db.WithContext(c).Scopes(opts...).Find(&clusters).Error
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) GetCntBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (count int64, err error) {
	err = db.WithContext(c).Scopes(opts...).Model(model.Cluster{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) DeleteBy(c context.Context, db *gorm.DB, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Delete(&model.Cluster{}).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) UpdateColumn(c context.Context, db *gorm.DB, name string, value interface{}, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Model(&model.Cluster{}).UpdateColumn(name, value).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}

func (s *cluster) Updates(c context.Context, db *gorm.DB, cluster *model.Cluster, opts ...store.DBOptions) (err error) {
	err = db.WithContext(c).Scopes(opts...).Updates(&cluster).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return
}
