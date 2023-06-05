package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/code-art/gin-im/util"
)

func Upload(c *gin.Context) {
	UploadToOSS(c)
}

// 上传到阿里云OSS
func UploadToOSS(c *gin.Context) {
	srcFile, head, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileType := c.PostForm("fileType")
	suffix := ".png"
	filename := head.Filename
	tem := strings.Split(filename, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	if fileType != "" {
		suffix = fileType
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，accessKeyId以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	endpoint := viper.GetString("oos.endpoint")
	accessKeyId := viper.GetString("oos.accessKeyId")
	accessKeySecret := viper.GetString("oos.accessKeySecret")
	bucket := viper.GetString("oos.bucket")

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bu, err := client.Bucket(bucket)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bu.PutObject(fileName, srcFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	accessAddress := viper.GetString("oos.accessAddress")
	url := fmt.Sprintf("%s/%s", accessAddress, fileName)
	c.JSON(http.StatusOK, util.R{
		Code: http.StatusOK,
		Msg:  "发送图片成功",
		Data: url,
	})
}

// 上传到本地
func UploadToLocal(c *gin.Context) {
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
