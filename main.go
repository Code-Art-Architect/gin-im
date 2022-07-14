package main

import (
	"github.com/code-art/gin-im/router"
	"github.com/code-art/gin-im/util"
)

func main() {
	util.InitConfig()
	util.InitMySQL()

	r := router.Router()
	r.Run(":8080")
}
