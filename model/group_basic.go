package model

import (
	"github.com/code-art/gin-im/util"
)

// 群信息
type GroupBasic struct {
	util.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (g *GroupBasic) TableName() string {
	return "group_basic"
}
