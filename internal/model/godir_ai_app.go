package model

type GodirAiApp struct {
	Base

	Name  string `gorm:"size:128;not null"`
	AppID string `gorm:"size:128;not null;index"`
	Desc  string `gorm:"size:500"`
	Icon  string `gorm:"size:255"`
	Cover string `gorm:"size:1000"`
}

func (GodirAiApp) TableName() string {
	return "godir_ai_app"
}
