package router

import (
	"github.com/gin-gonic/gin"
	"neuronet/internal/neuronetserver/controller/v1/permission"
	"neuronet/internal/neuronetserver/controller/v1/role"
	"neuronet/internal/neuronetserver/controller/v1/user"
	"neuronet/internal/neuronetserver/options"
	"neuronet/internal/pkg/message"
	"time"
)

func New(opts *options.Options) error {
	web := opts.WebServer.Get()
	versionGroup := web.Group("/v1")
	versionGroup.Use(opts.Interceptors.RequestID())

	healthGroup := versionGroup.Group("/health")
	{

		healthGroup.GET("", func(ctx *gin.Context) {
			type health struct {
				TIME string `json:"time"`
			}
			var myHealth = health{}
			myHealth.TIME = time.Now().Format("2006-01-02 15:04:05")
			message.WriteResponse(ctx, nil, myHealth)
		})
	}

	userController := user.NewController(opts.Db, opts.StoreFactory)
	userGroup := versionGroup.Group("/user")
	{
		userGroup.GET("/list", opts.Interceptors.JWTAuth(), userController.GetList)
		userGroup.POST("/register", opts.Interceptors.JWTAuth(), userController.Register)
		userGroup.POST("/login", userController.Login)
		userGroup.GET("", opts.Interceptors.JWTAuth(), userController.GetDetail)
		userGroup.DELETE("", opts.Interceptors.JWTAuth(), userController.Delete)
		userGroup.PUT("", opts.Interceptors.JWTAuth(), userController.Update)
	}

	roleController := role.NewController(opts.Db, opts.StoreFactory)
	roleGroup := versionGroup.Group("/user/role")
	{
		roleGroup.GET("/list", opts.Interceptors.JWTAuth(), roleController.GetList)
		roleGroup.POST("", opts.Interceptors.JWTAuth(), roleController.Create)
		roleGroup.DELETE("", opts.Interceptors.JWTAuth(), roleController.Delete)
		roleGroup.PUT("", opts.Interceptors.JWTAuth(), roleController.Update)
	}

	permissionController := permission.NewController(opts.Db, opts.StoreFactory)
	permissionGroup := versionGroup.Group("/user/permission")
	{
		permissionGroup.GET("/list", opts.Interceptors.JWTAuth(), permissionController.GetList)
		permissionGroup.POST("", opts.Interceptors.JWTAuth(), permissionController.Create)
		permissionGroup.DELETE("", opts.Interceptors.JWTAuth(), permissionController.Delete)
		permissionGroup.PUT("", opts.Interceptors.JWTAuth(), permissionController.Update)
	}

	return nil
}
