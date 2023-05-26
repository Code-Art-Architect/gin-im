package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} Welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	// htm, err := template.ParseFiles("index.html")
	// if err != nil {
	// 	panic(err)
	// }
	// htm.Execute(c.Writer, "index")
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!",
	// })
	c.HTML(http.StatusOK, "/chat/index.shtml", gin.H{
		"title": "Gin Framework",
	})
}
