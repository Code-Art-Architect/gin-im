package model

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/code-art/gin-im/util"
)

type Contact struct {
	gorm.Model
	OwnerId  uint   // 谁的关系信息
	TargetId uint   // 对应的谁
	Type     int    // 对应的类型 0 1 3
	Desc     string // 描述
}

func (c *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	var contacts []Contact
	var objIds []uint64
	var users []UserBasic
	util.DB.Where("owner_id = ?", userId).Find(&contacts)
	for _, contact := range contacts {
		fmt.Println(contact)
		objIds = append(objIds, uint64(contact.TargetId))
	}
	util.DB.Where("id in ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetId uint) int {
	if targetId != 0 {
		targetUser := FindUserById(targetId)
		if targetUser.Identity != "" {
			contact := Contact{
				OwnerId:  userId,
				TargetId: targetId,
				Type:     2,
			}
			util.DB.Create(&contact)
			return 1
		}
		return -1
	}
	return -1
}
