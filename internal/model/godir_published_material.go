package model

import (
	"gorm.io/gorm"
)

type GodirPublishedMaterial struct {
	gorm.Model

	UserID      uint   `gorm:"not null;index"` // User who published the material
	MaterialID  uint   `gorm:"not null;index"` // Reference to the material
	Description string `gorm:"type:text"`      // Description of the published material

	Control
}

func (GodirPublishedMaterial) TableName() string {
	return "godir_published_material"
}