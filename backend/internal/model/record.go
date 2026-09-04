package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CookRecord struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	DishID          uint64         `gorm:"not null;index" json:"dish_id"`
	CookedAt        time.Time      `gorm:"not null" json:"cooked_at"`
	Result          string         `gorm:"size:30;not null;default:normal" json:"result"`
	Rating          *float64       `gorm:"type:numeric(2,1)" json:"rating"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Changes         string         `gorm:"type:text" json:"changes"`
	FailureReason   string         `gorm:"type:text" json:"failure_reason"`
	NextImprovement string         `gorm:"type:text" json:"next_improvement"`
	Extra           datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"extra"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Photos []Attachment `gorm:"-" json:"photos,omitempty"`
}

func (CookRecord) TableName() string { return "cook_records" }
