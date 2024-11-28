package main

import (
	"github.com/zhufuyi/sponge/pkg/app"
	"gorm.io/gorm"
	initial2 "novaro-server/cmd/initial"
	"novaro-server/model"
)

var (
	secret = []byte("secret")
)

func main() {
	initial2.InitApp()
	services := initial2.CreateServices()
	closes := initial2.Close(services)
	//dropTableGorm()
	initDB()
	a := app.New(services, closes)
	a.Run()
}

func initDB() {
	db := model.GetDB()
	db.AutoMigrate(&model.Users{})
	db.AutoMigrate(&model.Collections{})
	db.AutoMigrate(&model.Comments{})
	db.AutoMigrate(&model.Events{})
	db.AutoMigrate(&model.Imgs{})
	db.AutoMigrate(&model.InvitationCodes{})
	db.AutoMigrate(&model.Invitations{})
	db.AutoMigrate(&model.NftInfo{})
	db.AutoMigrate(&model.PointsHistory{})
	db.AutoMigrate(&model.Posts{})
	db.AutoMigrate(&model.RePosts{})
	db.AutoMigrate(&model.Tags{})
	db.AutoMigrate(&model.TagsRecords{})
	db.AutoMigrate(&model.TwitterUsers{})
	db.AutoMigrate(&model.TwitterUserInfo{})
	db.AutoMigrate(&model.Users{})
	db.AutoMigrate(&model.PostPoints{})
	db.AutoMigrate(&model.NftLevel{})
	db.AutoMigrate(&model.PointsChangeLog{})
	db.AutoMigrate(&model.NftTokens{})
	db.AutoMigrate(&model.Likes{})
}
func dropTableGorm() {
	db := model.GetDB()

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("PRAGMA foreign_keys = OFF")

	db.Migrator().DropTable(&model.Posts{})

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("PRAGMA foreign_keys = ON")
}
