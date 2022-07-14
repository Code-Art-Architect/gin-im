package service

import (
	"github.com/code-art/gin-im/model"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 获取用户列表
// @Success 200 {string} json{"code", "message"}
// @Router /user/list [GET]
func GetUserList(c *gin.Context) {
	data := model.GetUserList()

	c.JSON(200, gin.H{
		"message": data,
	})
}
