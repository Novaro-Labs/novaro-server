package main

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"novaro-server/config"
	_ "novaro-server/docs"
	routes "novaro-server/router"
	"os/signal"
	"syscall"
)

var (
	secret = []byte("secret")
)

func main() {

	err := config.Init()
	if err != nil {
		log.Err(err).Msg("config init error")
	}
	defer config.Close()
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *graceful.Graceful {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	router, err := graceful.Default()
	router.Static("/assets", "./assets/")

	if err != nil {
		panic(err)
	}

	defer router.Close()
	router.Use(logger.SetLogger())
	router.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	e, err := casbin.NewEnforcer()
	if err != nil {
		panic(err)
	}
	//swag
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1", authz.NewAuthorizer(e))

	routes.AddHomeRoutes(v1)
	routes.AddAuthRoutes(v1)
	routes.AddOtherRoutes(v1)

	if err := router.RunWithContext(ctx); err != nil && err != context.Canceled {
		panic(err)
	}
	return router
}
