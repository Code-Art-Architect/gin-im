package service

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/code-art/gin-im/model"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/list [GET]
func GetUserList(c *gin.Context) {
	data := model.GetUserList()

	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 添加用户
// @Tags 用户模块
// @Param username query string false "用户名"
// @Param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/insert [POST]
func CreateUser(c *gin.Context) {
	user := model.UserBasic{}
	user.Name = c.Query("username")
	user.Password = c.Query("password")

	model.CreateUser(user)
	c.JSON(200, gin.H{
		"message": "添加用户成功！",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id query string false "id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/delete [DELETE]
func DeleteUser(c *gin.Context) {
	user := model.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)

	model.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功！",
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @Param id formData string false "id"
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Param phone formData string false "phone"
// @Param email formData string false "email"
// @Success 200 {string} json{"code", "message"}
// @Router /user/update [PUT]
func UpdateUser(c *gin.Context) {
	user := model.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"message": "修改参数不匹配！",
		})
	} else {
		model.UpdateUser(user)
		c.JSON(200, gin.H{
			"message": "更新用户成功！",
		})
	}
}
