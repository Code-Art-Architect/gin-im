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

func AddFriend(userId uint, targetId uint) (int, string) {
	if targetId != 0 {
		targetUser := FindUserById(targetId)
		if targetUser.Name != "" {
			tx := util.DB.Begin()
			defer tx.Commit()
			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
				}
			}()

			contact0 := Contact{}
			util.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetId).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不可以重复添加好友"
			}

			contact := Contact{
				OwnerId:  userId,
				TargetId: targetId,
				Type:     1,
			}
			util.DB.Create(&contact)

			contact1 := Contact{
				OwnerId:  targetId,
				TargetId: userId,
				Type:     1,
			}
			util.DB.Create(&contact1)

			return 1, ""
		}
		return -1, "添加失败"
	}
	return -1, "目标不存在"
}
