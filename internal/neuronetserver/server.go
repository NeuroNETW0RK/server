package neuronetserver

import (
	"NeuroNET/internal/neuronetserver/options"
	"NeuroNET/internal/neuronetserver/router"
	"NeuroNET/pkg/log"
	"NeuroNET/pkg/shutdown"
	"NeuroNET/pkg/shutdown/shutdownmanagers/posixsignal"
	"os"
	"os/signal"
	"syscall"
)

type cloudServer struct {
	gs      *shutdown.GracefulShutdown
	options *options.Options
}

func New(configName string) *cloudServer {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	opts, err := options.NewOptions(configName)
	if err != nil {
		log.Fatalf("unable to initialize command options: %v", err)
	}
	err = opts.Complete()
	if err != nil {
		log.Fatalf("unable to complete options: %v", err)
	}
	return &cloudServer{
		gs:      gs,
		options: opts,
	}
}

func (d *cloudServer) Run() {
	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	defer log.Sync()
	// 初始化 APIs 路由
	err := router.New(d.options)
	if err != nil {
		log.Fatalf("failed to init routers: %v ", err)
	}
	// 启动服务
	d.runServer()
	<-stopCh
}

func (d *cloudServer) runServer() {
	d.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		db, _ := d.options.Db.DB()
		db.Close()
		d.options.WebServer.Close()
		return nil
	}))

	// start shutdown managers
	if err := d.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	d.options.WebServer.Run()
}
