package role

import (
	v1 "NeuroNET/internal/neuronetserver/dto/v1"
	"NeuroNET/internal/neuronetserver/service/v1/role"
	"NeuroNET/internal/neuronetserver/store"
	"NeuroNET/internal/pkg/code"
	"NeuroNET/internal/pkg/message"
	"NeuroNET/pkg/errors"
	"NeuroNET/pkg/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		roleService: role.NewService(db, store),
	}
}

type controller struct {
	roleService role.Service
}

func (c *controller) Create(ctx *gin.Context) {
	log.C(ctx).Infof("create role")

	args := new(v1.RoleCreateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.roleService.Create(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("role get list function called")

	args := new(v1.RoleListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.roleService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("role delete function called")

	args := new(v1.RoleDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.roleService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) Update(ctx *gin.Context) {
	log.C(ctx).Infof("role update function called")

	args := new(v1.RoleUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.roleService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
