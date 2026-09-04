package model

import (
	"time"

	"gorm.io/datatypes"
)

type AppSetting struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	SettingKey   string         `gorm:"size:100;not null;unique" json:"setting_key"`
	SettingValue datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"setting_value"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (AppSetting) TableName() string { return "app_settings" }
