package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

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
