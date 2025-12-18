package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"size:128;not null"`
	Password string `gorm:"size:255;not null"` // 存储加密后的密码
	Source   int64  `gorm:""`
	Control
}

func (User) TableName() string {
	return "user"
}
