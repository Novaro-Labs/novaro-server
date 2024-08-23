package config

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var RDB *redis.Client

var ClientId string
var ClientSecret string
var Proxy string

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	var err error

	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	ClientId = viper.GetString("client_id")
	ClientSecret = viper.GetString("client_secret")
	Proxy = viper.GetString("proxy")

	env := viper.GetString("env")
	if env != "dev" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			"127.0.0.1",
			"postgres",
			"sqltest",
			"novaro",
			5432,
		)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Database connection test failed:", err)
		}
	} else {
		DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
