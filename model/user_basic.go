package model

import (
	"fmt"
	"time"

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
	Salt          string 
	LoginTime     time.Time
	HeartBeatTime time.Time
	LogOutTime    time.Time
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

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	util.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	util.DB.Where("name = ? and password = ?", name, password).First(&user)
	
	// token加密
	temp := util.Md5Encode(fmt.Sprintf("%d", time.Now().Unix()))
	util.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	user.Identity = temp
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return util.DB.Where("phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return util.DB.Where("email = ?", email).First(&user)
}
