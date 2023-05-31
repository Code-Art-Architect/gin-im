package test_gorm

import (
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/code-art/gin-im/model"
	"github.com/code-art/gin-im/util"
)

var db *gorm.DB

func init() {
	var dsn = "root:root1234@tcp(127.0.0.1:3306)/gin-im?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to Connect DataBase")
	}
	util.InitConfig()
	util.InitMySQL()
}

func TestMigrate(t *testing.T) {
	// _ = db.AutoMigrate(&model.UserBasic{})
	// _ = db.AutoMigrate(&model.Message{})
	// _ = db.AutoMigrate(&model.Contact{})
	// _ = db.AutoMigrate(&model.GroupBasic{})
	_ = db.AutoMigrate(&model.Community{})
}

func TestAddUser(t *testing.T) {
	util.DB.Create(&model.UserBasic{
		Name:          "特朗普",
		Password:      util.Md5Encode("123456"),
		Phone:         "13770367889",
		Email:         "152944@gmail.com",
		DeviceInfo:    "Macbook Pro M1 MAX",
		LoginTime:     time.Now(),
		LogOutTime:    time.Now(),
		HeartBeatTime: time.Now(),
	})
}
