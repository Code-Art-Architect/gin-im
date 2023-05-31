package router

import (
	"net/http"

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

	// 静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("view/**/*.html")

	// 首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)

	// 好友关系
	r.GET("/contact/load-friends", service.SearchFriend)
	r.POST("/contact/add-friend", service.AddFriend)
	r.POST("/contact/create-community", service.CreateCommunity)

	// 通用页面跳转
	r.GET("/:path1/:path2.shtml", func(c *gin.Context) {
		url := c.Request.URL.Path
		c.HTML(http.StatusOK, url, nil)
	})

	// 用户模块
	r.GET("/user/list", service.GetUserList)
	r.POST("/user/insert", service.CreateUser)
	r.DELETE("/user/delete", service.DeleteUser)
	r.PUT("/user/update", service.UpdateUser)
	r.POST("/user/login", service.Login)
	r.GET("/user/find", service.FindUser)

	// 发送消息
	r.GET("/user/msg", service.SendMessage)
	r.GET("/user/sendUMsg", service.SendUserMessage)

	// 上传文件
	r.POST("/attach/upload", service.Upload)

	return r
}
