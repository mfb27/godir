package model

import (
	"gorm.io/gorm"
)

// GodirPublishedLike represents a like by a user for a published material
type GodirPublishedLike struct {
	gorm.Model

	UserID      uint `gorm:"not null"`
	PublishedID uint `gorm:"not null"`

	Control
}

func (GodirPublishedLike) TableName() string {
	return "godir_published_like"
}
