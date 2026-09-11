package repository

import (
	"testing"

	"food-memory/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupDishList(t *testing.T) *DishRepo {
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
	if err := db.AutoMigrate(&model.Dish{}, &model.Tag{}, &model.DishTag{}); err != nil {
		t.Fatal(err)
	}
	return NewDishRepo(db)
}

func seedListedDish(t *testing.T, r *DishRepo, name string, tags []model.Tag) *model.Dish {
	t.Helper()
	dish := &model.Dish{Name: name, Status: model.StatusWantToCook, Extra: datatypes.JSON([]byte("{}"))}
	if err := r.Create(dish); err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		return dish
	}
	for i := range tags {
		if err := r.db.Create(&tags[i]).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := r.ReplaceTags(nil, dish.ID, tags); err != nil {
		t.Fatal(err)
	}
	return dish
}

func TestListUntaggedOnlyHasDishesWithoutDisplayTags(t *testing.T) {
	r := setupDishList(t)
	rice := seedListedDish(t, r, "红烧肉", []model.Tag{{Name: "下饭", Type: model.TagCustom}})
	pork := seedListedDish(t, r, "排骨汤", []model.Tag{{Name: "排骨", Type: model.TagIngredient}})
	plain := seedListedDish(t, r, "白粥", nil)

	untagged := true
	items, total, err := r.List(model.DishListQuery{Untagged: &untagged, Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("total=%d want 2", total)
	}
	got := map[string]bool{}
	for _, d := range items {
		got[d.Name] = true
	}
	if !got[pork.Name] || !got[plain.Name] {
		t.Fatalf("got %#v, want 排骨汤 and 白粥", got)
	}
	if got[rice.Name] {
		t.Fatal("带展示标签的菜不应出现在其他")
	}
}
