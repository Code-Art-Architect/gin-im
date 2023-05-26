package test_gorm

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/code-art/gin-im/model"
)

var db *gorm.DB

func init() {
	var dsn = "root:root1234@tcp(127.0.0.1:3306)/gin-im?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to Connect DataBase")
	}
}

func TestMigrate(t *testing.T) {
	// _ = db.AutoMigrate(&model.UserBasic{})
	// _ = db.AutoMigrate(&model.Message{})
	_ = db.AutoMigrate(&model.Contact{})
	_ = db.AutoMigrate(&model.GroupBasic{})
}
