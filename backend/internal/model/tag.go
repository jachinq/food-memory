package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:50;not null;uniqueIndex:uk_tags_name_type" json:"name"`
	Type       string         `gorm:"size:30;not null;default:custom;uniqueIndex:uk_tags_name_type" json:"type"`
	UsageCount int            `gorm:"not null;default:0" json:"usage_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tag) TableName() string { return "tags" }
