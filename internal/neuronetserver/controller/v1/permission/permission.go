package permission

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/service/v1/permission"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/message"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
)

var _ Controller = (*controller)(nil)

type Controller interface {
	Create(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
}

func NewController(db *gorm.DB, store store.Factory) *controller {
	return &controller{
		permissionService: permission.NewService(db, store),
	}
}

type controller struct {
	permissionService permission.Service
}

func (c *controller) Create(ctx *gin.Context) {
	log.C(ctx).Infof("create permission")

	args := new(v1.PermissionCreateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.permissionService.Create(ctx, args)
	if err != nil {
		message.FailedWithMsg(ctx, err.Error(), err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("permission get list function called")

	args := new(v1.PermissionListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.permissionService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("permission delete function called")

	args := new(v1.PermissionDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.permissionService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) Update(ctx *gin.Context) {
	log.C(ctx).Infof("permission update function called")

	args := new(v1.PermissionUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.permissionService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
