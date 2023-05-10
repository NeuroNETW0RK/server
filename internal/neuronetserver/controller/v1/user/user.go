package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/service/v1/user"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/message"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
)

var _ Controller = (*controller)(nil)

type Controller interface {
	Register(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Login(ctx *gin.Context)
	Update(ctx *gin.Context)
}

func NewController(db *gorm.DB, store store.Factory) *controller {
	return &controller{
		userService: user.NewService(db, store),
	}
}

type controller struct {
	userService user.Service
}

func (c *controller) Register(ctx *gin.Context) {
	log.C(ctx).Infof("user register function called")

	args := new(v1.UserRegisterArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.userService.Register(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) GetList(ctx *gin.Context) {
	log.C(ctx).Infof("user get list function called")

	args := new(v1.UserListArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.userService.GetList(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) GetDetail(ctx *gin.Context) {
	log.C(ctx).Infof("user get detail function called")

	args := new(v1.UserDetailArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.userService.GetDetail(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Delete(ctx *gin.Context) {
	log.C(ctx).Infof("user delete function called")

	args := new(v1.UserDeleteArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.userService.Delete(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}

func (c *controller) Login(ctx *gin.Context) {
	log.C(ctx).Infof("user login function called")

	args := new(v1.UserLoginArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	reply, err := c.userService.Login(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, reply)
}

func (c *controller) Update(ctx *gin.Context) {
	log.C(ctx).Infof("user update function called")

	args := new(v1.UserUpdateArgs)

	if err := ctx.ShouldBind(&args); err != nil {
		message.Failed(ctx, errors.WithCode(code.ErrBind, err.Error()))
		return
	}

	err := c.userService.Update(ctx, args)
	if err != nil {
		message.Failed(ctx, err)
		return
	}

	message.Success(ctx, nil)
}
