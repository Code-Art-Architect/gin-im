package model

import (
	"github.com/code-art/gin-im/util"
)

type Community struct {
	util.Model
	Name    string `json:"name,omitempty"`
	OwnerId uint   `json:"ownerId,omitempty"`
	Icon    string `json:"icon,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

func (c Community) TableName() string {
	return "community"
}

func CreateCommunity(c Community) (int, string) {
	if len(c.Name) == 0 {
		return -1, "群名称太短"
	}
	if c.OwnerId == 0 {
		return -1, "请先登录"
	}
	if err := util.DB.Create(&c).Error; err != nil {
		return -1, "建群失败"
	}
	return 1, "建群成功"
}

func LoadCommunity(userId uint) []*Community {
	var data []*Community
	var contacts []*Contact
	var groupIds []uint
	util.DB.Where("owner_id = ? and type = 2", userId).Find(&contacts)
	for _, c := range contacts {
		groupIds = append(groupIds, c.TargetId)
	}

	util.DB.Where("id in ?", groupIds).Find(&data)
	return data
}