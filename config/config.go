package config

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB                      *gorm.DB
	RDB                     *redis.Client
	RBQ                     *amqp091.Connection
	InvitatioCodeExpiration time.Duration
	InvitatioCodeLength     int
	ClientId                string
	ClientSecret            string
	Proxy                   string
	UploadPath 				string
)


func Init() error {
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
	UploadPath = viper.GetString("uploadPath")
	// 初始化数据库连接
	err = initDB()

	// 初始化 Redis 连接
	err = initRedis()

	// 初始化 RabbitMQ 连接
	err = initRabbitMQ()

	return err
}

func initDB() error {
	env := viper.GetString("env")
	var err error

	if env != "dev" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			viper.GetString("database.host"),
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("database.name"),
			viper.GetInt("database.port"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	return nil
}

func initRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 可以添加一个 Ping 操作来验证连接
	_, err := RDB.Ping(RDB.Context()).Result()
	if err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	return nil
}

func initRabbitMQ() error {
	var err error
	RBQ, err = amqp091.Dial("amqp://root:rabbitmqtest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("RabbitMQ 连接失败: %w", err)
	}
	return nil
}

// Close 关闭所有连接
func Close() {
	if DB != nil {
		db, _ := DB.DB()
		db.Close()
	}
	if RDB != nil {
		RDB.Close()
	}
	if RBQ != nil {
		RBQ.Close()
	}
}
