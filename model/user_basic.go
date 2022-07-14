package model

import (
	"fmt"

	"github.com/code-art/gin-im/util"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIP      string
	ClientPort    string
	LoginTime     uint64
	HeartBeatTime uint64
	LogOutTime    uint64
	IsLogOut      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	util.DB.Find(&data)

	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	return util.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return util.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return util.DB.Model(&user).Updates(
		UserBasic{
			Name:     user.Name,
			Password: user.Password,
			Phone:    user.Phone,
			Email:    user.Email,
		})
}
