package main

import (
	"context"
	"novaro-server/config"
	_ "novaro-server/docs"
	"novaro-server/model"
	"novaro-server/src/routes"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	secret = []byte("secret")
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	defer config.Close()

	setConfig()

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

func setConfig() {
	// 迁移数据库
	DB = config.DB
	DB.AutoMigrate(&model.Collections{})
	DB.AutoMigrate(&model.Comments{})
	DB.AutoMigrate(&model.Posts{})
	DB.AutoMigrate(&model.RePosts{})
	DB.AutoMigrate(&model.Tags{})
	DB.AutoMigrate(&model.TagsRecords{})
	DB.AutoMigrate(&model.Users{})
	DB.AutoMigrate(&model.TwitterUsers{})
	DB.AutoMigrate(&model.InvitationCodes{})
	DB.AutoMigrate(&model.Invitations{})

	// 创建 cron 实例
	c := cron.New()
	// 添加定时任务：每分钟执行同步
	c.AddFunc("@every 5m", func() {
		model.SyncToDatabase()
	})

	c.AddFunc("@every 3s", func() {
		model.SyncCountToDataBase()
	})
	c.Start()
}
