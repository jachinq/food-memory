package repository

import (
	"food-memory/internal/model"

	"gorm.io/gorm"
)

type AttachmentRepo struct {
	db *gorm.DB
}

func NewAttachmentRepo(db *gorm.DB) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

func (r *AttachmentRepo) Create(att *model.Attachment) error {
	return r.db.Create(att).Error
}

func (r *AttachmentRepo) Bind(ids []uint64, bizType string, bizID uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&model.Attachment{}).
		Where("id IN ?", ids).
		Updates(map[string]any{"biz_type": bizType, "biz_id": bizID}).Error
}

func (r *AttachmentRepo) BindByURLs(urls []string, bizType string, bizID uint64) error {
	if len(urls) == 0 {
		return nil
	}
	return r.db.Model(&model.Attachment{}).
		Where("file_url IN ?", urls).
		Updates(map[string]any{"biz_type": bizType, "biz_id": bizID}).Error
}

func (r *AttachmentRepo) ListByBiz(bizType string, bizID uint64) ([]model.Attachment, error) {
	var items []model.Attachment
	err := r.db.Where("biz_type = ? AND biz_id = ?", bizType, bizID).
		Order("sort_order ASC, id ASC").
		Find(&items).Error
	return items, err
}

func (r *AttachmentRepo) ListByBizIDs(bizType string, bizIDs []uint64) ([]model.Attachment, error) {
	if len(bizIDs) == 0 {
		return nil, nil
	}
	var items []model.Attachment
	err := r.db.Where("biz_type = ? AND biz_id IN ?", bizType, bizIDs).
		Order("sort_order ASC, id ASC").
		Find(&items).Error
	return items, err
}

func (r *AttachmentRepo) GetByID(id uint64) (*model.Attachment, error) {
	var att model.Attachment
	if err := r.db.First(&att, id).Error; err != nil {
		return nil, err
	}
	return &att, nil
}
