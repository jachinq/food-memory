package service

import (
	"encoding/json"
	"testing"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupDishService(t *testing.T) *DishService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
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
	if err := db.AutoMigrate(&model.Dish{}, &model.CookRecord{}, &model.Tag{}, &model.DishTag{}, &model.Attachment{}); err != nil {
		t.Fatal(err)
	}
	return NewDishService(db, repository.NewDishRepo(db), repository.NewTagRepo(db), repository.NewAttachmentRepo(db), repository.NewRecordRepo(db))
}

func TestUpdateKeepsCustomTagsOnReload(t *testing.T) {
	s := setupDishService(t)
	created, err := s.Create(model.DishCreateInput{Name: "红烧肉", Status: model.StatusWantToCook})
	if err != nil {
		t.Fatal(err)
	}

	payload := []byte(`{"name":"红烧肉","status":"want_to_cook","tags":[{"name":"下饭","type":"custom"}]}`)
	var in model.DishCreateInput
	if err := json.Unmarshal(payload, &in); err != nil {
		t.Fatal(err)
	}
	updated, err := s.Update(created.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	if !hasTag(updated.Tags, "下饭", model.TagCustom) {
		t.Fatalf("update response tags=%v, want custom 下饭", names(updated.Tags))
	}

	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !hasTag(got.Tags, "下饭", model.TagCustom) {
		t.Fatalf("reload tags=%v, want custom 下饭 (edit page would show empty)", names(got.Tags))
	}
}

func names(tags []model.Tag) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.Name+"/"+t.Type)
	}
	return out
}

func hasTag(tags []model.Tag, name, typ string) bool {
	for _, t := range tags {
		if t.Name == name && t.Type == typ {
			return true
		}
	}
	return false
}
