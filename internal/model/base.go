package model

type Control struct {
	CreatedBy int64 `gorm:""`
	UpdatedBy int64 `gorm:""`
}
