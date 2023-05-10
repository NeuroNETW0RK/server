package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
	"time"
)

type WebServerOptions struct {
	Mode    string
	Address string
	Port    int
}

type WebServer struct {
	engine  *gin.Engine
	Port    int
	Address string
	server  *http.Server
}

func (web *WebServer) Get() *gin.Engine {
	return web.engine
}

func (web *WebServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	web.engine.ServeHTTP(w, req)
}

func (web *WebServer) Run() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", web.Address, web.Port),
		Handler: web.engine,
	}
	web.server = srv

	go func() {
		log.Infof("Start to listening the incoming requests on http address: %v", web.Port)

		if err := web.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(err.Error())
		}

		log.Infof("Server on %s stopped", web.Port)

	}()
}

func (web *WebServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := web.server.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}
}

func NewWebServer(opt *WebServerOptions) (*WebServer, error) {

	gin.SetMode(opt.Mode)
	server := WebServer{engine: gin.Default(), Address: opt.Address, Port: opt.Port}

	server.engine.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	}))

	return &server, nil
}
