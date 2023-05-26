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
	c.HTML(http.StatusOK, "/chat/index.shtml", gin.H{
		"title": "Gin Framework",
	})
}
