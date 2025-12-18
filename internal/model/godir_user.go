package model

import "gorm.io/gorm"

type GodirUser struct {
	gorm.Model

	Username string `gorm:"size:128;not null"`
	Password string `gorm:"size:255;not null"` // 存储加密后的密码
	
	// 新增用户信息字段
	Avatar  string `gorm:"size:255"` // 头像URL
	Nickname string `gorm:"size:128"` // 昵称
	Gender   int    `gorm:"default:0"` // 性别: 0-未知, 1-男, 2-女

	Control
}

func (GodirUser) TableName() string {
	return "godir_user"
}