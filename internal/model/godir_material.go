package model

import "gorm.io/gorm"

type GodirMaterial struct {
	Model[GodirMaterial]

	Base
	UserID           uint   `gorm:"not null;index"`
	FileName         string `gorm:"size:255;not null"`
	FileExt          string `gorm:"size:10;not null"`
	FileSize         int64  `gorm:"not null"`
	ContentType      string `gorm:"size:100"`
	OssBucket        string `gorm:"size:100;not null"`
	OssFilePath      string `gorm:"size:500;not null"`
	CoverOssFilePath string `gorm:"size:500"` // Cover 为封面/缩略图信息
	Control

	// CoverURL string `gorm:"size:1000"`
	// URL      string `gorm:"size:1000"`
}

func (GodirMaterial) TableName() string {
	return "godir_material"
}

func NewGodirMaterial(db *gorm.DB) *GodirMaterial {
	m := &GodirMaterial{}
	m.db = db
	return m
}
