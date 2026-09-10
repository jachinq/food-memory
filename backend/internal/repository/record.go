package repository

import (
	"food-memory/internal/model"
	"time"

	"gorm.io/gorm"
)

type RecordRepo struct {
	db *gorm.DB
}

func NewRecordRepo(db *gorm.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (r *RecordRepo) Create(tx *gorm.DB, rec *model.CookRecord) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(rec).Error
}

func (r *RecordRepo) Update(rec *model.CookRecord) error {
	return r.db.Save(rec).Error
}

func (r *RecordRepo) GetByID(id uint64) (*model.CookRecord, error) {
	var rec model.CookRecord
	if err := r.db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecordRepo) ListByDish(dishID uint64) ([]model.CookRecord, error) {
	var items []model.CookRecord
	err := r.db.Where("dish_id = ?", dishID).Order("cooked_at DESC, id DESC").Find(&items).Error
	return items, err
}

func (r *RecordRepo) SoftDelete(id uint64) error {
	return r.db.Delete(&model.CookRecord{}, id).Error
}

func (r *RecordRepo) Stats(tx *gorm.DB, dishID uint64) (count int64, last *time.Time, avg *float64, err error) {
	if tx == nil {
		tx = r.db
	}
	if err = tx.Model(&model.CookRecord{}).Where("dish_id = ?", dishID).Count(&count).Error; err != nil {
		return
	}
	if count == 0 {
		return
	}
	var latest model.CookRecord
	if err = tx.Where("dish_id = ?", dishID).Order("cooked_at DESC, id DESC").First(&latest).Error; err != nil {
		return
	}
	last = &latest.CookedAt
	var avgVal *float64
	if err = tx.Model(&model.CookRecord{}).
		Where("dish_id = ? AND rating IS NOT NULL", dishID).
		Select("AVG(rating)").
		Scan(&avgVal).Error; err != nil {
		return
	}
	avg = avgVal
	return
}
