package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/code-art/gin-im/doc"
	"github.com/code-art/gin-im/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/insert", service.CreateUser)
	r.DELETE("/user/delete", service.DeleteUser)
	r.PUT("/user/update", service.UpdateUser)
	r.POST("/user/login", service.Login)

	// 发送消息
	r.GET("/user/msg", service.SendMessage)
	r.GET("/user/sendUMsg", service.SendUserMessage)

	return r
}
