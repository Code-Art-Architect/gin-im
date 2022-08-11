package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *gorm.DB
	Redigo *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app: ", viper.Get("app"))
	fmt.Println("config mysql: ", viper.Get("mysql"))
}

func InitMySQL() {
	// 自定义日志模块打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢sql的阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	dns := viper.GetString("mysql.dns")
	DB, _ = gorm.Open(mysql.Open(dns), &gorm.Config{Logger: newLogger})
}

func InitRedis() {
	Redigo = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConnection"),
	})
}

const (
	PublishKey = "webSocket"
)

// 发送消息到Redis
func PublishToRedis(c context.Context, channel, msg string) error {
	var err error
	fmt.Println("Published: ", msg)

	err = Redigo.Publish(c, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}

	return err
}

// 订阅Redis
func SubscribeFromRedis(c context.Context, channel string) (string, error) {
	sub := Redigo.Subscribe(c, channel)
	fmt.Println("Subscribed: ", sub)
	m, err := sub.ReceiveMessage(c)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("Subscribed: ", m.Payload)
	return m.Payload, err
}
