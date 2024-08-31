package config

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var RDB *redis.Client

var InvitatioCodeExpiration time.Duration
var InvitatioCodeLength int
var ClientId string
var ClientSecret string
var Proxy string

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("invitation_code_expire_day", "7")
	viper.SetDefault("invitation_code_length", "8")
	viper.AddConfigPath(".")
	var err error

	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	expire := viper.GetInt("invitation_code_expire_day")
	InvitatioCodeExpiration = time.Duration(float64(expire) * 24 * float64(time.Hour))
	InvitatioCodeLength = viper.GetInt("invitation_code_length")
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
