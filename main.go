package main

import (
	"time"

	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/router"
	"github.com/code-art/gin-im/util"
)

func InitTimer() {
	util.Timer(time.Second*3, time.Second*6, model.ClearConnection, "")
}

func main() {
	util.InitConfig()
	util.InitMySQL()
	util.InitRedis()
	InitTimer()

	r := router.Router()
	_ = r.Run(":8080")
}
