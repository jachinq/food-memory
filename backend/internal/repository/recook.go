package repository

import (
	"food-memory/internal/model"
	"time"

	"gorm.io/gorm"
)

type RecookRepo struct {
	db *gorm.DB
}

func NewRecookRepo(db *gorm.DB) *RecookRepo {
	return &RecookRepo{db: db}
}

func (r *RecookRepo) List(status string) ([]model.RecookPlan, error) {
	q := r.db.Preload("Dish").Preload("Dish.Tags").Model(&model.RecookPlan{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var items []model.RecookPlan
	err := q.Order("planned_date ASC NULLS LAST, created_at DESC").Find(&items).Error
	return items, err
}

func (r *RecookRepo) Create(plan *model.RecookPlan) error {
	return r.db.Create(plan).Error
}

func (r *RecookRepo) GetByID(id uint64) (*model.RecookPlan, error) {
	var plan model.RecookPlan
	if err := r.db.Preload("Dish").First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *RecookRepo) Update(plan *model.RecookPlan) error {
	return r.db.Save(plan).Error
}

func (r *RecookRepo) SoftDelete(id uint64) error {
	return r.db.Delete(&model.RecookPlan{}, id).Error
}

func (r *RecookRepo) ActiveByDish(dishID uint64) (*model.RecookPlan, error) {
	var plan model.RecookPlan
	err := r.db.Where("dish_id = ? AND status = ?", dishID, model.RecookActive).
		Order("id DESC").First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *RecookRepo) Complete(id uint64) error {
	now := time.Now()
	return r.db.Model(&model.RecookPlan{}).Where("id = ?", id).Updates(map[string]any{
		"status":       model.RecookCompleted,
		"completed_at": now,
	}).Error
}
