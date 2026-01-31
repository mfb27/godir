package model

// GodirPublishedLike represents a like by a user for a published material
type GodirPublishedLike struct {
	Base

	UserID      uint `gorm:"not null"`
	PublishedID uint `gorm:"not null"`
}

func (GodirPublishedLike) TableName() string {
	return "godir_published_like"
}
