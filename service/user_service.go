package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// Login
// @Summary 登录
// @Tags 用户模块
// @Param username query string false "用户名"
// @Param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/login [POST]
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	user := model.FindUserByName(username)
	if user.Name == "" {
		c.JSON(400, gin.H{
			"message": "该用户不存在",
		})
		return
	}

	flag := util.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(400, gin.H{
			"message": "密码错误",
		})
		return
	}
	
	pwd := util.MakePassword(password, user.Salt)
	data := model.FindUserByNameAndPwd(username, pwd) 

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
	user := model.UserBasic{
		Name: c.Query("username"),
		Password: c.Query("password"),
	}
	user.LoginTime = time.Now()
	user.LogOutTime = time.Now()
	user.HeartBeatTime = time.Now()
	
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	
	data := model.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(400, gin.H{
			"message": "用户名已经存在！",
		})
		return
	}
	
	user.Password = util.MakePassword(user.Password, salt)

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

// 防止跨域站点伪造请求
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println()
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := util.SubscribeFromRedis(c, util.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息：", msg)

		t := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s][%s]", t, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}
