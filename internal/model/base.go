package model

import "gorm.io/gorm"

type Control struct {
	CreatedBy int64 `gorm:""`
	UpdatedBy int64 `gorm:""`
}

type Base struct {
	gorm.Model
	Control
}
