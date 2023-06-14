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

var _ TagController = (*tagController)(nil)

type TagController interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Info(ctx *gin.Context)
}

func NewTagController(db *gorm.DB, store store.Factory) *tagController {
	return &tagController{
		imageTagService: image.NewTagService(db, store),
	}
}

type tagController struct {
	imageTagService image.TagService
}

func (c *tagController) Create(ctx *gin.Context) {
	log.C(ctx).Infof("create image function called")

	args := new(v1.ImageCreateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.imageTagService.Create(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *tagController) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("delete image function called")

	args := new(v1.ImageDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.imageTagService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *tagController) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("get image list function called")

	args := new(v1.ImageListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.imageTagService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *tagController) Info(ctx *gin.Context) {
	log.C(ctx).Infof("get image list function called")

	args := new(v1.ImageListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.imageTagService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *tagController) Update(ctx *gin.Context) {
	log.C(ctx).Infof("update image function called")

	args := new(v1.ImageUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.imageTagService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
