package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/util"
)

// 创建群聊
func CreateCommunity(c *gin.Context) {
	var community model.Community
	if err := c.ShouldBindJSON(&community); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	i, s := model.CreateCommunity(community)
	if i == -1 {
		c.JSON(http.StatusBadRequest, util.R{
			Code: http.StatusBadRequest,
			Msg:  s,
		})
		return
	}

	c.JSON(http.StatusOK, util.R{
		Code: http.StatusOK,
		Msg:  s,
	})
}

// 加载群列表
func LoadCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	data := model.LoadCommunity(uint(userId))
	c.JSON(http.StatusOK, util.R{
		Code: http.StatusOK,
		Data: data,
	})
}
