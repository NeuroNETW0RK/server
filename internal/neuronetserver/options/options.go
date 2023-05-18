package options

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/configs"
	"neuronet/internal/neuronetserver/router/interceptor"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/neuronetserver/store/mysql"
	"neuronet/internal/pkg/core"
	"neuronet/pkg/db"
	"neuronet/pkg/k8s"
	"neuronet/pkg/log"
	cache "neuronet/pkg/redis"
)

type Options struct {
	ComponentConfig configs.Config
	StoreFactory    store.Factory
	WebServer       *core.WebServer
	Interceptors    interceptor.Interceptor
	Db              *gorm.DB
	Redis           *redis.Client
}

func NewOptions(configName string) (*Options, error) {
	configs.InitConfig(configName)
	return &Options{
		ComponentConfig: configs.Get(),
	}, nil
}

// Complete completes all the required options
func (o *Options) Complete() error {
	if err := o.register(); err != nil {
		return err
	}
	return nil
}

func (o *Options) register() error {
	// 注册日志
	if err := o.registerLogger(); err != nil {
		return err
	}
	fmt.Println("[REGISTER] register logger successful")
	// 注册数据库
	if err := o.registerDatabase(); err != nil {
		return err
	}
	fmt.Println("[REGISTER] register db successful")
	// 注册中间件
	if err := o.registerInterceptor(); err != nil {
		return err
	}
	fmt.Println("[REGISTER] register middleware successful")
	// 注册web服务器
	if err := o.registerWebServer(); err != nil {
		return err
	}
	fmt.Println("[REGISTER] register web server successful")
	// 注册clusterSets
	if err := o.registerK8s(); err != nil {
		return err
	}
	fmt.Println("[REGISTER] register k8s successful")
	return nil
}

func (o *Options) registerLogger() error {
	logConf := o.ComponentConfig.Log
	logOptions := &log.Options{
		DisableCaller:     logConf.DisableCaller,
		DisableStacktrace: logConf.DisableStacktrace,
		Level:             logConf.Level,
		Format:            logConf.Format,
		OutputPaths:       logConf.OutputPaths,
	}
	log.Init(logOptions)
	return nil
}

func (o *Options) registerDatabase() (err error) {
	dbConf := o.ComponentConfig.Mysql
	dbOptions := &db.MySQLOptions{
		Host:                  dbConf.Addr,
		Username:              dbConf.User,
		Password:              dbConf.Password,
		Database:              dbConf.Database,
		MaxIdleConnections:    dbConf.MaxIdleConn,
		MaxOpenConnections:    dbConf.MaxOpenConn,
		MaxConnectionLifeTime: dbConf.ConnMaxLifeTime,
		LogLevel:              dbConf.LogLevel,
	}

	o.Db, err = db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	o.StoreFactory = mysql.NewMysqlDatastore()

	return nil
}

func (o *Options) registerInterceptor() error {
	o.Interceptors = interceptor.New()
	return nil
}

func (o *Options) registerWebServer() error {
	webServerConfig := o.ComponentConfig.Server
	webServerOptions := &core.WebServerOptions{
		Mode: webServerConfig.Mode,
		Port: webServerConfig.Port,
	}
	webServer, err := core.NewWebServer(
		webServerOptions,
	)
	if err != nil {
		return err
	}
	o.WebServer = webServer
	return nil
}

func (o *Options) registerRedis() error {
	redisConf := o.ComponentConfig.Redis
	redisOpt := &cache.Options{
		Addr:       redisConf.Addr,
		Port:       redisConf.Port,
		Db:         redisConf.Db,
		MaxReTries: redisConf.MaxReTries,
		Password:   redisConf.Password,
		PoolSize:   redisConf.PoolSize,
	}
	client, err := cache.New(redisOpt)
	if err != nil {
		return err
	}
	o.Redis = client
	return nil
}

func (o *Options) registerK8s() error {
	k8s.NewClusterSet()
	k8s.NewCoreV1Store(k8s.GetClusterSets())
	ctx := context.Background()
	clusterBos, err := o.StoreFactory.Cluster().GetListBy(ctx, o.Db)
	if err != nil {
		return err
	}
	clusterSets := k8s.GetClusterSets()
	for _, clusterBo := range clusterBos {
		stop := make(chan struct{})
		clientSets, err := k8s.NewClientSets(ctx, clusterBo.ConfigPath, stop)
		if err != nil {
			close(stop)
			log.Warnf("[REGISTER] load cluster %s failed, err: %v", clusterBo.Name, err)
			fmt.Printf("[REGISTER] load cluster %s failed, err: %v \n", clusterBo.Name, err)
			continue
		}
		clusterSets.Add(clusterBo.Name, clientSets)
	}
	return nil
}
