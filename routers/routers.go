package routers

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "novaro-server/docs"
	"novaro-server/middlewares"
	"os/signal"
	"syscall"
)

var (
	secret = []byte("secret")
)

// NewRouter create a new router
func NewRouter() *graceful.Graceful {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r, err := graceful.Default()
	r.Static("/assets", "./assets/")

	if err != nil {
		panic(err)
	}

	defer r.Close()
	r.Use(logger.SetLogger())
	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	e, err := casbin.NewEnforcer()
	if err != nil {
		panic(err)
	}
	//swag
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middlewares.Cors())
	v1 := r.Group("/v1", authz.NewAuthorizer(e))

	AddHomeRoutes(v1)
	AddAuthRoutes(v1)
	AddOtherRoutes(v1)

	if err := r.RunWithContext(ctx); err != nil && err != context.Canceled {
		panic(err)
	}
	return r
}
