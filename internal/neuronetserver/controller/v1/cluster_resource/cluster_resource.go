package clusterresource

import (
	"github.com/gin-gonic/gin"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	clusterresource "neuronet/internal/neuronetserver/service/v1/cluster_resouce"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/message"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
)

type Controller interface {
	SingleNodes(ctx *gin.Context)
	GroupNodes(ctx *gin.Context)
	AllNodes(ctx *gin.Context)
}

func NewController() *controller {
	return &controller{
		clusterResourceService: clusterresource.NewService(),
	}
}

type controller struct {
	clusterResourceService clusterresource.Service
}

func (c *controller) SingleNodes(ctx *gin.Context) {
	log.C(ctx).Infof("get single node function called")

	args := new(v1.SingleNodeArgs)

	if err := ctx.ShouldBindUri(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.clusterResourceService.SingleNodes(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) GroupNodes(ctx *gin.Context) {
	log.C(ctx).Infof("get group node function called")

	args := new(v1.GroupNodesArgs)

	if err := ctx.ShouldBindUri(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.clusterResourceService.GroupNodes(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) AllNodes(ctx *gin.Context) {
	log.C(ctx).Infof("get all node function called")

	args := new(v1.AllNodesArgs)

	if err := ctx.ShouldBindUri(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.clusterResourceService.AllNodes(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}
