package model

import "gorm.io/gorm"

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
