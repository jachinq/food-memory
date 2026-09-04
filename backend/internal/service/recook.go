package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"gorm.io/gorm"
)

type RecookService struct {
	db     *gorm.DB
	plans  *repository.RecookRepo
	dishes *repository.DishRepo
}

func NewRecookService(db *gorm.DB, plans *repository.RecookRepo, dishes *repository.DishRepo) *RecookService {
	return &RecookService{db: db, plans: plans, dishes: dishes}
}

func (s *RecookService) List(status string) ([]model.RecookPlan, error) {
	return s.plans.List(status)
}

func (s *RecookService) Create(in model.RecookPlanInput) (*model.RecookPlan, error) {
	if in.DishID == 0 {
		return nil, fmt.Errorf("请选择菜品")
	}
	if _, err := s.dishes.GetByID(in.DishID); err != nil {
		return nil, err
	}
	plan := &model.RecookPlan{
		DishID: in.DishID,
		Status: model.RecookActive,
		Reason: strings.TrimSpace(in.Reason),
	}
	if d := strings.TrimSpace(in.PlannedDate); d != "" {
		t, err := time.ParseInLocation("2006-01-02", d, time.Local)
		if err != nil {
			return nil, fmt.Errorf("计划日期格式无效")
		}
		plan.PlannedDate = &t
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(plan).Error; err != nil {
			return err
		}
		return tx.Model(&model.Dish{}).Where("id = ?", in.DishID).
			Update("status", model.StatusWantToRecook).Error
	})
	if err != nil {
		return nil, err
	}
	return s.plans.GetByID(plan.ID)
}

func (s *RecookService) Update(id uint64, in model.RecookPlanUpdateInput) (*model.RecookPlan, error) {
	plan, err := s.plans.GetByID(id)
	if err != nil {
		return nil, err
	}
	if d := strings.TrimSpace(in.PlannedDate); d != "" {
		t, err := time.ParseInLocation("2006-01-02", d, time.Local)
		if err != nil {
			return nil, fmt.Errorf("计划日期格式无效")
		}
		plan.PlannedDate = &t
	}
	if in.Reason != "" {
		plan.Reason = in.Reason
	}
	if in.Status != "" {
		plan.Status = in.Status
	}
	if err := s.plans.Update(plan); err != nil {
		return nil, err
	}
	return s.plans.GetByID(id)
}

func (s *RecookService) Complete(id uint64) (*model.RecookPlan, error) {
	if _, err := s.plans.GetByID(id); err != nil {
		return nil, err
	}
	if err := s.plans.Complete(id); err != nil {
		return nil, err
	}
	return s.plans.GetByID(id)
}

func (s *RecookService) Cancel(id uint64) error {
	plan, err := s.plans.GetByID(id)
	if err != nil {
		return err
	}
	plan.Status = model.RecookCancelled
	return s.plans.Update(plan)
}

func (s *RecookService) Delete(id uint64) error {
	if _, err := s.plans.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}
	return s.plans.SoftDelete(id)
}
