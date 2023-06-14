package router

import (
	"neuronet/internal/neuronetserver/controller/v1/cluster"
	clusterresource "neuronet/internal/neuronetserver/controller/v1/cluster_resource"
	"neuronet/internal/neuronetserver/controller/v1/image"
	"neuronet/internal/neuronetserver/controller/v1/permission"
	"neuronet/internal/neuronetserver/controller/v1/role"
	"neuronet/internal/neuronetserver/controller/v1/user"
	"neuronet/internal/neuronetserver/options"
	"neuronet/internal/pkg/message"
	"time"

	"github.com/gin-gonic/gin"
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

	clusterController := cluster.NewController(opts.Db, opts.StoreFactory)
	clusterGroup := versionGroup.Group("/cluster")
	{
		clusterGroup.GET("/list", clusterController.GetList)
		clusterGroup.POST("", clusterController.Create)
		clusterGroup.DELETE("", clusterController.Delete)
		clusterGroup.PUT("", clusterController.Update)
		clusterGroup.POST("/reload", clusterController.Reload)
	}

	clusterResourceController := clusterresource.NewController()
	clusterResourceGroup := versionGroup.Group("/cluster/:cluster_name")
	{
		clusterResourceGroup.GET("/:node_name", clusterResourceController.SingleNodes)
		clusterResourceGroup.GET("/label/:label", clusterResourceController.GroupNodes)
		clusterResourceGroup.GET("", clusterResourceController.AllNodes)
	}

	imageController := image.NewController(opts.Db, opts.StoreFactory)
	imageGroup := versionGroup.Group("/image")
	{
		// image
		imageGroup.GET("/list", imageController.GetList)
		imageGroup.GET("/info/:image_id", imageController.Info)
		imageGroup.POST("", imageController.Create)
		imageGroup.DELETE("/:image_id", clusterController.Delete)
		imageGroup.PUT("/:image_id", clusterController.Update)
		// imageTag
		tagGroup := imageGroup.Group("/tag")
		{
			tagGroup.GET("/list", imageController.GetList)
			imageGroup.GET("/info/:image_tag_id", imageController.Info)
			imageGroup.POST("", imageController.Create)
			imageGroup.DELETE("/:image_tag_id", clusterController.Delete)
			imageGroup.PUT("/:image_tag_id", clusterController.Update)
		}
		// imageBuild
		buildGroup := imageGroup.Group("/build")
		{
			buildGroup.GET("/list", imageController.GetList)
			buildGroup.GET("/info/:image_build_id", imageController.Info)
			buildGroup.POST("", imageController.Create)
			buildGroup.DELETE("/:image_build_id", clusterController.Delete)
			buildGroup.PUT("/:image_build_id", clusterController.Update)
		}
	}
	return nil
}
