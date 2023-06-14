package image

import (
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/service/v1/image"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/message"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ BuildController = (*buildController)(nil)

type BuildController interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Info(ctx *gin.Context)
}

func NewBuildController(db *gorm.DB, store store.Factory) *buildController {
	return &buildController{
		buildService: image.NewBuildService(db, store),
	}
}

type buildController struct {
	buildService image.BuildService
}

func (c *buildController) Create(ctx *gin.Context) {
	log.C(ctx).Infof("create image function called")

	args := new(v1.ImageCreateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.buildService.Create(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *buildController) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("delete image function called")

	args := new(v1.ImageDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.buildService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *buildController) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("get image list function called")

	args := new(v1.ImageListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.buildService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *buildController) Info(ctx *gin.Context) {
	log.C(ctx).Infof("get image list function called")

	args := new(v1.ImageListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.buildService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *buildController) Update(ctx *gin.Context) {
	log.C(ctx).Infof("update image function called")

	args := new(v1.ImageUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.buildService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
