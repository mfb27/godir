package model

import (
	"gorm.io/gorm"
)

type GodirMaterial struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	FileName    string `gorm:"size:255;not null"`
	FileSize    int64  `gorm:"not null"`
	ContentType string `gorm:"size:100"`
	Bucket      string `gorm:"size:100;not null"`
	Key         string `gorm:"size:500;not null"`
	URL         string `gorm:"size:1000"`
	// Cover 为封面/缩略图信息
	CoverKey string `gorm:"size:500"` // 存储在对象存储的 key
	CoverURL string `gorm:"size:1000"`
	Control
}

func (GodirMaterial) TableName() string {
	return "godir_material"
}