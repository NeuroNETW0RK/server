package cluster

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/service/v1/cluster"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/message"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
)

var _ Controller = (*controller)(nil)

type Controller interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Reload(ctx *gin.Context)
}

func NewController(db *gorm.DB, store store.Factory) *controller {
	return &controller{
		clusterService: cluster.NewService(db, store),
	}
}

type controller struct {
	clusterService cluster.Service
}

func (c *controller) Create(ctx *gin.Context) {
	log.C(ctx).Infof("create cluster function called")

	args := new(v1.ClusterCreateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.clusterService.Create(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("delete cluster function called")

	args := new(v1.ClusterDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.clusterService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("get cluster list function called")

	args := new(v1.ClusterListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.clusterService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Update(ctx *gin.Context) {
	log.C(ctx).Infof("update cluster function called")

	args := new(v1.ClusterUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.clusterService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) Reload(ctx *gin.Context) {
	log.C(ctx).Infof("reload cluster function called")

	args := new(v1.ClusterReloadArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.clusterService.Reload(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
