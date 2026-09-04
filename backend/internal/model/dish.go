package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Dish struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:120;not null" json:"name"`
	CoverImageURL    string         `gorm:"type:text" json:"cover_image_url"`
	SourceURL        string         `gorm:"type:text" json:"source_url"`
	SourcePlatform   string         `gorm:"size:50" json:"source_platform"`
	Description      string         `gorm:"type:text" json:"description"`
	Status           string         `gorm:"size:30;not null;default:want_to_cook" json:"status"`
	Rating           *float64       `gorm:"type:numeric(2,1)" json:"rating"`
	Difficulty       *int16         `json:"difficulty"`
	CookTimeMinutes  *int           `json:"cook_time_minutes"`
	MainIngredients  string         `gorm:"type:text" json:"main_ingredients"`
	Taste            string         `gorm:"size:80" json:"taste"`
	Scene            string         `gorm:"size:80" json:"scene"`
	Note             string         `gorm:"type:text" json:"note"`
	LastCookedAt     *time.Time     `json:"last_cooked_at"`
	CookCount        int            `gorm:"not null;default:0" json:"cook_count"`
	IsFavorite       bool           `gorm:"not null;default:false" json:"is_favorite"`
	Extra            datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"extra"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	CoverThumbnailURL string        `gorm:"-" json:"cover_thumbnail_url"`
	Tags              []Tag         `gorm:"many2many:dish_tags;" json:"tags,omitempty"`
	Records           []CookRecord  `gorm:"foreignKey:DishID" json:"records,omitempty"`
}

func (Dish) TableName() string { return "dishes" }

type DishTag struct {
	DishID    uint64    `gorm:"primaryKey" json:"dish_id"`
	TagID     uint64    `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (DishTag) TableName() string { return "dish_tags" }
