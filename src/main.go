package main

import (
	"context"
	"novaro-server/config"
	"novaro-server/model"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"

	"novaro-server/src/routes"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var secret = []byte("secret")

func init() {
	config.Init()
	db := config.DB
	db.AutoMigrate(&model.Collections{})
	db.AutoMigrate(&model.Comments{})
	db.AutoMigrate(&model.Posts{})
	db.AutoMigrate(&model.RePosts{})
	db.AutoMigrate(&model.Tags{})
	db.AutoMigrate(&model.TagsRecords{})
	db.AutoMigrate(&model.Users{})
	db.AutoMigrate(&model.TwitterUsers{})
	db.AutoMigrate(&model.InvitationCodes{})
	db.AutoMigrate(&model.Invitations{})

	// 创建 cron 实例
	c := cron.New()
	// 添加定时任务：每分钟执行同步
	c.AddFunc("@every 1m", func() {
		model.SyncData()
	})
	c.Start()
}

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

	router.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	e, err := casbin.NewEnforcer()
	if err != nil {
		panic(err)
	}
	v1 := router.Group("/v1", authz.NewAuthorizer(e))

	routes.AddHomeRoutes(v1)
	routes.AddAuthRoutes(v1)
	routes.AddOtherRoutes(v1)

	if err := router.RunWithContext(ctx); err != nil && err != context.Canceled {
		panic(err)
	}
	return router
}
