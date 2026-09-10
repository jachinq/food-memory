package service

import (
	"fmt"
	"strings"
	"time"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecordService struct {
	db     *gorm.DB
	recs   *repository.RecordRepo
	dishes *repository.DishRepo
	atts   *repository.AttachmentRepo
	plans  *repository.RecookRepo
}

func NewRecordService(db *gorm.DB, recs *repository.RecordRepo, dishes *repository.DishRepo, atts *repository.AttachmentRepo, plans *repository.RecookRepo) *RecordService {
	return &RecordService{db: db, recs: recs, dishes: dishes, atts: atts, plans: plans}
}

func (s *RecordService) List(dishID uint64) ([]model.CookRecord, error) {
	if _, err := s.dishes.GetByID(dishID); err != nil {
		return nil, err
	}
	items, err := s.recs.ListByDish(dishID)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	atts, err := s.atts.ListByBizIDs(model.BizCookRecord, ids)
	if err != nil {
		return nil, err
	}
	grouped := map[uint64][]model.Attachment{}
	for _, a := range atts {
		grouped[a.BizID] = append(grouped[a.BizID], a)
	}
	for i := range items {
		items[i].Photos = grouped[items[i].ID]
	}
	return items, nil
}

func (s *RecordService) Create(dishID uint64, in model.RecordInput) (*model.RecordSaveResult, error) {
	if _, err := s.dishes.GetByID(dishID); err != nil {
		return nil, err
	}
	rec, err := s.buildRecord(dishID, in)
	if err != nil {
		return nil, err
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.recs.Create(tx, rec); err != nil {
			return err
		}
		if err := s.bindPhotos(tx, rec.ID, in); err != nil {
			return err
		}
		if err := s.recalcDish(tx, dishID, in.UpdateDishStatus); err != nil {
			return err
		}
		if err := s.plans.CompleteActiveByDish(tx, dishID); err != nil {
			return err
		}
		if in.UpdateDishStatus == model.StatusWantToRecook {
			plan := &model.RecookPlan{DishID: dishID, Status: model.RecookActive}
			if err := tx.Create(plan).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	saved, err := s.recs.GetByID(rec.ID)
	if err != nil {
		return nil, err
	}
	photos, _ := s.atts.ListByBiz(model.BizCookRecord, saved.ID)
	saved.Photos = photos
	return &model.RecordSaveResult{Record: *saved, SuggestStatus: suggestStatus(in)}, nil
}

func (s *RecordService) Update(id uint64, in model.RecordInput) (*model.RecordSaveResult, error) {
	rec, err := s.recs.GetByID(id)
	if err != nil {
		return nil, err
	}
	built, err := s.buildRecord(rec.DishID, in)
	if err != nil {
		return nil, err
	}
	rec.CookedAt = built.CookedAt
	rec.Result = built.Result
	rec.Rating = built.Rating
	rec.Notes = built.Notes
	rec.Changes = built.Changes
	rec.FailureReason = built.FailureReason
	rec.NextImprovement = built.NextImprovement

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(rec).Error; err != nil {
			return err
		}
		if err := s.bindPhotos(tx, rec.ID, in); err != nil {
			return err
		}
		return s.recalcDishAfterEdit(tx, rec.DishID)
	})
	if err != nil {
		return nil, err
	}
	photos, _ := s.atts.ListByBiz(model.BizCookRecord, rec.ID)
	rec.Photos = photos
	return &model.RecordSaveResult{Record: *rec, SuggestStatus: suggestStatus(in)}, nil
}

func (s *RecordService) Delete(id uint64) error {
	rec, err := s.recs.GetByID(id)
	if err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.CookRecord{}, id).Error; err != nil {
			return err
		}
		return s.recalcDishAfterEdit(tx, rec.DishID)
	})
}

func (s *RecordService) buildRecord(dishID uint64, in model.RecordInput) (*model.CookRecord, error) {
	result := in.Result
	if result == "" {
		result = model.ResultNormal
	}
	if !model.ValidResult(result) {
		return nil, fmt.Errorf("无效的制作结果")
	}
	cookedAt := time.Now()
	if strings.TrimSpace(in.CookedAt) != "" {
		parsed, err := parseTime(in.CookedAt)
		if err != nil {
			return nil, fmt.Errorf("制作日期格式无效")
		}
		cookedAt = parsed
	}
	return &model.CookRecord{
		DishID:          dishID,
		CookedAt:        cookedAt,
		Result:          result,
		Rating:          in.Rating,
		Notes:           in.Notes,
		Changes:         in.Changes,
		FailureReason:   in.FailureReason,
		NextImprovement: in.NextImprovement,
		Extra:           datatypes.JSON([]byte("{}")),
	}, nil
}

func (s *RecordService) bindPhotos(tx *gorm.DB, recordID uint64, in model.RecordInput) error {
	repo := repository.NewAttachmentRepo(tx)
	if err := repo.UnbindByBiz(model.BizCookRecord, recordID); err != nil {
		return err
	}
	if err := repo.Bind(in.AttachmentIDs, model.BizCookRecord, recordID); err != nil {
		return err
	}
	return repo.BindByURLs(in.PhotoURLs, model.BizCookRecord, recordID)
}

func (s *RecordService) recalcDish(tx *gorm.DB, dishID uint64, status string) error {
	count, last, avg, err := s.recs.Stats(tx, dishID)
	if err != nil {
		return err
	}
	updates := map[string]any{
		"cook_count":     count,
		"last_cooked_at": last,
	}
	if avg != nil {
		updates["rating"] = round1(*avg)
	} else {
		updates["rating"] = nil
	}
	if status != "" {
		if !model.ValidDishStatus(status) {
			return fmt.Errorf("无效的菜品状态")
		}
		updates["status"] = status
	}
	return tx.Model(&model.Dish{}).Where("id = ?", dishID).Updates(updates).Error
}

func (s *RecordService) recalcDishAfterEdit(tx *gorm.DB, dishID uint64) error {
	if err := s.recalcDish(tx, dishID, ""); err != nil {
		return err
	}
	var n int64
	if err := tx.Model(&model.RecookPlan{}).
		Where("dish_id = ? AND status = ?", dishID, model.RecookActive).
		Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	status, err := dishStatusFromLatestCook(tx, dishID)
	if err != nil {
		return err
	}
	return tx.Model(&model.Dish{}).Where("id = ?", dishID).Update("status", status).Error
}

func suggestStatus(in model.RecordInput) string {
	if in.Result == model.ResultSuccess {
		return model.StatusSuccess
	}
	if in.Rating != nil && *in.Rating >= 4 {
		return model.StatusSuccess
	}
	return ""
}

func parseTime(v string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	var last error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, v, time.Local)
		if err == nil {
			return t, nil
		}
		last = err
	}
	return time.Time{}, last
}
