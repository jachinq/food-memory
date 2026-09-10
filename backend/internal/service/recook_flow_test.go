package service

import (
	"errors"
	"testing"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type recookHarness struct {
	db      *gorm.DB
	recook  *RecookService
	records *RecordService
	plans   *repository.RecookRepo
	dishes  *repository.DishRepo
}

func setupRecook(t *testing.T) *recookHarness {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.Dish{}, &model.CookRecord{}, &model.RecookPlan{}, &model.Attachment{}); err != nil {
		t.Fatal(err)
	}
	dishRepo := repository.NewDishRepo(db)
	recordRepo := repository.NewRecordRepo(db)
	attRepo := repository.NewAttachmentRepo(db)
	planRepo := repository.NewRecookRepo(db)
	recook := NewRecookService(db, planRepo, dishRepo)
	records := NewRecordService(db, recordRepo, dishRepo, attRepo, planRepo)
	return &recookHarness{db: db, recook: recook, records: records, plans: planRepo, dishes: dishRepo}
}

func seedDish(t *testing.T, h *recookHarness, status string) *model.Dish {
	t.Helper()
	dish := &model.Dish{Name: "红烧肉", Status: status, Extra: datatypes.JSON([]byte("{}"))}
	if err := h.dishes.Create(dish); err != nil {
		t.Fatal(err)
	}
	return dish
}

func TestCreateCookRecordCompletesActivePlan(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.records.Create(dish.ID, model.RecordInput{Result: model.ResultSuccess}); err != nil {
		t.Fatal(err)
	}
	got, err := h.plans.GetByID(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.RecookCompleted {
		t.Fatalf("status=%s want completed", got.Status)
	}
	if got.CompletedAt == nil {
		t.Fatal("expected completed_at")
	}
}

func TestFailedCookRecordCompletesActivePlan(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusWantToRecook)
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.records.Create(dish.ID, model.RecordInput{Result: model.ResultFailed}); err != nil {
		t.Fatal(err)
	}
	got, err := h.plans.GetByID(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.RecookCompleted {
		t.Fatalf("status=%s want completed", got.Status)
	}
}

func TestEmptyCompleteRejected(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.recook.Complete(plan.ID); !errors.Is(err, ErrEmptyRecookComplete) {
		t.Fatalf("err=%v want ErrEmptyRecookComplete", err)
	}
	got, err := h.plans.GetByID(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.RecookActive {
		t.Fatalf("status=%s want active", got.Status)
	}
}

func TestSecondActivePlanRejected(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	if _, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID}); err != nil {
		t.Fatal(err)
	}
	if _, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID}); !errors.Is(err, ErrRecookAlreadyActive) {
		t.Fatalf("err=%v want ErrRecookAlreadyActive", err)
	}
}

func TestCancelRecalcsDishStatusFromLatestRecord(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	if _, err := h.records.Create(dish.ID, model.RecordInput{Result: model.ResultFailed}); err != nil {
		t.Fatal(err)
	}
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := h.recook.Cancel(plan.ID); err != nil {
		t.Fatal(err)
	}
	got, err := h.dishes.GetByID(dish.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.StatusFailed {
		t.Fatalf("status=%s want failed", got.Status)
	}
}

func TestCancelWithoutRecordsSetsWantToCook(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := h.recook.Cancel(plan.ID); err != nil {
		t.Fatal(err)
	}
	got, err := h.dishes.GetByID(dish.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.StatusWantToCook {
		t.Fatalf("status=%s want want_to_cook", got.Status)
	}
}

func TestWantToRecookOnNewRecordOpensNewPlan(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	old, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.records.Create(dish.ID, model.RecordInput{
		Result:           model.ResultSuccess,
		UpdateDishStatus: model.StatusWantToRecook,
	}); err != nil {
		t.Fatal(err)
	}
	gotOld, err := h.plans.GetByID(old.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotOld.Status != model.RecookCompleted {
		t.Fatalf("old status=%s want completed", gotOld.Status)
	}
	active, err := h.plans.ActiveByDish(dish.ID)
	if err != nil {
		t.Fatal(err)
	}
	if active.ID == old.ID {
		t.Fatal("expected a new active plan")
	}
	dishGot, err := h.dishes.GetByID(dish.ID)
	if err != nil {
		t.Fatal(err)
	}
	if dishGot.Status != model.StatusWantToRecook {
		t.Fatalf("dish status=%s want want_to_recook", dishGot.Status)
	}
}

func TestUpdateRecordDoesNotCompletePlan(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	saved, err := h.records.Create(dish.ID, model.RecordInput{Result: model.ResultNormal})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := h.records.Update(saved.Record.ID, model.RecordInput{Result: model.ResultSuccess}); err != nil {
		t.Fatal(err)
	}
	got, err := h.plans.GetByID(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.RecookActive {
		t.Fatalf("status=%s want active", got.Status)
	}
}

func TestDeleteRecordDoesNotReopenCompletedPlan(t *testing.T) {
	h := setupRecook(t)
	dish := seedDish(t, h, model.StatusSuccess)
	plan, err := h.recook.Create(model.RecookPlanInput{DishID: dish.ID})
	if err != nil {
		t.Fatal(err)
	}
	saved, err := h.records.Create(dish.ID, model.RecordInput{Result: model.ResultSuccess})
	if err != nil {
		t.Fatal(err)
	}
	if err := h.records.Delete(saved.Record.ID); err != nil {
		t.Fatal(err)
	}
	got, err := h.plans.GetByID(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != model.RecookCompleted {
		t.Fatalf("status=%s want completed", got.Status)
	}
}
