package main

import (
	"time"

	"github.com/spf13/viper"

	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/router"
	"github.com/code-art/gin-im/util"
)

func InitTimer() {
	delayHeartbeat := viper.GetInt64("task.delayHeartbeat")
	heartbeatHz := viper.GetInt("task.heartbeatHz")
	util.Timer(time.Duration(delayHeartbeat)*time.Second, time.Duration(heartbeatHz)*time.Second, model.ClearConnection, "")
}

func main() {
	util.InitConfig()
	util.InitMySQL()
	util.InitRedis()
	model.InitUDP()
	InitTimer()

	var httpPort = viper.GetString("server.port.http")
	r := router.Router()
	_ = r.Run(":" + httpPort)
}
