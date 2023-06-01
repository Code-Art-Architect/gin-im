package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/util"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} Welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := model.UserBasic{
		Model:    util.Model{ID: uint(userId)},
		Identity: token,
	}
	c.HTML(http.StatusOK, "/chat/index.shtml", user)
}
