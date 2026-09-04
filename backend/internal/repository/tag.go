package repository

import (
	"food-memory/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) List(tagType, keyword string) ([]model.Tag, error) {
	q := r.db.Model(&model.Tag{})
	if tagType != "" {
		q = q.Where("type = ?", tagType)
	}
	if keyword != "" {
		q = q.Where("name ILIKE ?", "%"+keyword+"%")
	}
	var items []model.Tag
	err := q.Order("usage_count DESC, id DESC").Find(&items).Error
	return items, err
}

func (r *TagRepo) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepo) Upsert(tx *gorm.DB, name, tagType string) (*model.Tag, error) {
	if tx == nil {
		tx = r.db
	}
	tag := model.Tag{Name: name, Type: tagType}
	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(&tag).Error
	if err != nil {
		return nil, err
	}
	if tag.ID == 0 {
		if err := tx.Where("name = ? AND type = ?", name, tagType).First(&tag).Error; err != nil {
			return nil, err
		}
	}
	return &tag, nil
}

func (r *TagRepo) RecalcUsage(tx *gorm.DB, tagIDs []uint64) error {
	if tx == nil {
		tx = r.db
	}
	if len(tagIDs) == 0 {
		return nil
	}
	return tx.Model(&model.Tag{}).Where("id IN ?", tagIDs).Updates(map[string]any{
		"usage_count": gorm.Expr("(SELECT COUNT(1) FROM dish_tags WHERE dish_tags.tag_id = tags.id)"),
		"updated_at":  gorm.Expr("NOW()"),
	}).Error
}
