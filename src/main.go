package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/novaro-server/src/routes"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *graceful.Graceful {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router, err := graceful.Default()
	if err != nil {
		panic(err)
	}
	defer router.Close()

	router.Use(logger.SetLogger())

	e, err := casbin.NewEnforcer()
	if err != nil {
		panic(err)
	}
	v1 := router.Group("/v1", authz.NewAuthorizer(e))

	routes.AddHomeRoutes(v1)
	routes.AddAuthRoutes(v1)

	if err := router.RunWithContext(ctx); err != nil && err != context.Canceled {
		panic(err)
	}
	return router
}
