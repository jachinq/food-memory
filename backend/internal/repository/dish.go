package repository

import (
	"food-memory/internal/model"

	"gorm.io/gorm"
)

type DishRepo struct {
	db *gorm.DB
}

func NewDishRepo(db *gorm.DB) *DishRepo {
	return &DishRepo{db: db}
}

func (r *DishRepo) DB() *gorm.DB { return r.db }

func (r *DishRepo) Create(dish *model.Dish) error {
	return r.db.Create(dish).Error
}

func (r *DishRepo) Update(dish *model.Dish) error {
	return r.db.Save(dish).Error
}

func (r *DishRepo) GetByID(id uint64) (*model.Dish, error) {
	var dish model.Dish
	err := r.db.Preload("Tags").First(&dish, id).Error
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (r *DishRepo) SoftDelete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("dish_id = ?", id).Delete(&model.CookRecord{}).Error; err != nil {
			return err
		}
		if err := tx.Where("dish_id = ?", id).Delete(&model.RecookPlan{}).Error; err != nil {
			return err
		}
		if err := tx.Where("dish_id = ?", id).Delete(&model.DishTag{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Dish{}, id).Error
	})
}

func (r *DishRepo) List(q model.DishListQuery) ([]model.Dish, int64, error) {
	db := r.db.Model(&model.Dish{})
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where(`
			name ILIKE ? OR main_ingredients ILIKE ?
			OR id IN (
				SELECT dt.dish_id FROM dish_tags dt
				JOIN tags t ON t.id = dt.tag_id AND t.deleted_at IS NULL
				WHERE t.name ILIKE ?
				AND t.type NOT IN (?, ?)
			)
		`, like, like, like, model.TagTaste, model.TagScene)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.Untagged != nil && *q.Untagged {
		db = db.Where(`id NOT IN (
			SELECT dt.dish_id FROM dish_tags dt
			JOIN tags t ON t.id = dt.tag_id AND t.deleted_at IS NULL
			WHERE t.type NOT IN (?, ?, ?)
		)`, model.TagIngredient, model.TagTaste, model.TagScene)
	} else if q.Tag != "" {
		tagLike := "%" + q.Tag + "%"
		db = db.Where(`id IN (
			SELECT dt.dish_id FROM dish_tags dt
			JOIN tags t ON t.id = dt.tag_id AND t.deleted_at IS NULL
			WHERE t.name ILIKE ? OR CAST(t.id AS TEXT) = ?
		)`, tagLike, q.Tag)
	}
	if q.Cooked != nil {
		if *q.Cooked {
			db = db.Where("cook_count > 0")
		} else {
			db = db.Where("cook_count = 0")
		}
	}
	if q.Success != nil {
		if *q.Success {
			db = db.Where("status = ?", model.StatusSuccess)
		} else {
			db = db.Where("status <> ?", model.StatusSuccess)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	var items []model.Dish
	err := db.Preload("Tags").
		Order("updated_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&items).Error
	return items, total, err
}

func (r *DishRepo) ReplaceTags(tx *gorm.DB, dishID uint64, tags []model.Tag) error {
	if tx == nil {
		tx = r.db
	}
	var dish model.Dish
	if err := tx.First(&dish, dishID).Error; err != nil {
		return err
	}
	return tx.Model(&dish).Association("Tags").Replace(tags)
}
