package model

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	BizType      string         `gorm:"size:50;not null;index:idx_attachments_biz" json:"biz_type"`
	BizID        uint64         `gorm:"not null;index:idx_attachments_biz" json:"biz_id"`
	FileName     string         `gorm:"size:255;not null" json:"file_name"`
	FileURL      string         `gorm:"type:text;not null" json:"file_url"`
	ThumbnailURL string         `gorm:"type:text" json:"thumbnail_url"`
	MimeType     string         `gorm:"size:80" json:"mime_type"`
	FileSize     int64          `json:"file_size"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	SortOrder    int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Attachment) TableName() string { return "attachments" }
