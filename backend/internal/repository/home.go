package repository

import (
	"food-memory/internal/model"
	"time"

	"gorm.io/gorm"
)

type HomeRepo struct {
	db *gorm.DB
}

func NewHomeRepo(db *gorm.DB) *HomeRepo {
	return &HomeRepo{db: db}
}

func (r *HomeRepo) Stats() (model.HomeStats, error) {
	var s model.HomeStats
	if err := r.db.Model(&model.Dish{}).Count(&s.TotalDishes).Error; err != nil {
		return s, err
	}
	if err := r.db.Model(&model.Dish{}).Where("cook_count > 0").Count(&s.CookedCount).Error; err != nil {
		return s, err
	}
	if err := r.db.Model(&model.Dish{}).Where("status = ?", model.StatusSuccess).Count(&s.SuccessCount).Error; err != nil {
		return s, err
	}
	start := time.Now().In(time.Local)
	monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	if err := r.db.Model(&model.CookRecord{}).Where("cooked_at >= ?", monthStart).Count(&s.MonthCookCount).Error; err != nil {
		return s, err
	}
	return s, nil
}

func (r *HomeRepo) Recent(limit int) ([]model.Dish, error) {
	var items []model.Dish
	err := r.db.Preload("Tags").
		Where("last_cooked_at IS NOT NULL").
		Order("last_cooked_at DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *HomeRepo) OverdueHighRating(limit int) ([]model.Dish, error) {
	cutoff := time.Now().Add(-30 * 24 * time.Hour)
	var items []model.Dish
	err := r.db.Preload("Tags").
		Where("rating >= ? AND status IN ? AND (last_cooked_at IS NULL OR last_cooked_at < ?)",
			4.0, []string{model.StatusSuccess, model.StatusWantToRecook}, cutoff).
		Order("last_cooked_at ASC NULLS FIRST").
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *HomeRepo) RandomOld() (*model.Dish, error) {
	var dish model.Dish
	err := r.db.Preload("Tags").
		Where("cook_count > 0 AND status IN ?", []string{model.StatusSuccess, model.StatusWantToRecook, model.StatusCooked}).
		Order("RANDOM()").
		First(&dish).Error
	if err != nil {
		return nil, err
	}
	return &dish, nil
}

func (r *HomeRepo) WantToCook(limit int) ([]model.Dish, error) {
	var items []model.Dish
	err := r.db.Preload("Tags").
		Where("status = ?", model.StatusWantToCook).
		Order("updated_at DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}
