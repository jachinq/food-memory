package model

import (
	"time"

	"gorm.io/gorm"
)

type RecookPlan struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	DishID      uint64         `gorm:"not null" json:"dish_id"`
	PlannedDate *time.Time     `gorm:"type:date" json:"planned_date"`
	Status      string         `gorm:"size:30;not null;default:active" json:"status"`
	Reason      string         `gorm:"type:text" json:"reason"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Dish *Dish `gorm:"foreignKey:DishID" json:"dish,omitempty"`
}

func (RecookPlan) TableName() string { return "recook_plans" }
