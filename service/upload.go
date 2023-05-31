package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
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

	fileType := c.PostForm("fileType")

	suffix := ".png"
	filename := srcFile.Filename
	tem := strings.Split(filename, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}

	if fileType != "" {
		suffix = fileType
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
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
