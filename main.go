package main

import (
	"github.com/code-art/gin-im/router"
	"github.com/code-art/gin-im/util"
)

func main() {
	util.InitConfig()
	util.InitMySQL()
	util.InitRedis()
	util.InitTimer()

	r := router.Router()
	_ = r.Run(":8080")
}
