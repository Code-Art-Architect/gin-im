package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/code-art/gin-im/util"
)

func Upload(c *gin.Context) {
	srcFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), ".png")
	url := "./asset/upload/" + fileName
	err = c.SaveUploadedFile(srcFile, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, util.R{
		Code: http.StatusOK,
		Msg:  "发送图片成功",
		Data: url,
	})
}
