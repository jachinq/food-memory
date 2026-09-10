package service

import (
	"errors"
	"fmt"
	"strings"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DishService struct {
	db    *gorm.DB
	dishes *repository.DishRepo
	tags  *repository.TagRepo
	atts  *repository.AttachmentRepo
	recs  *repository.RecordRepo
}

func NewDishService(db *gorm.DB, dishes *repository.DishRepo, tags *repository.TagRepo, atts *repository.AttachmentRepo, recs *repository.RecordRepo) *DishService {
	return &DishService{db: db, dishes: dishes, tags: tags, atts: atts, recs: recs}
}

func (s *DishService) Create(in model.DishCreateInput) (*model.Dish, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, fmt.Errorf("菜名为必填")
	}
	status := in.Status
	if status == "" {
		status = model.StatusWantToCook
	}
	if !model.ValidDishStatus(status) {
		return nil, fmt.Errorf("无效的菜品状态")
	}
	dish := &model.Dish{
		Name:            name,
		CoverImageURL:   strings.TrimSpace(in.CoverImageURL),
		SourceURL:       strings.TrimSpace(in.SourceURL),
		SourcePlatform:  strings.TrimSpace(in.SourcePlatform),
		Description:     strings.TrimSpace(in.Description),
		Status:          status,
		Rating:          in.Rating,
		Difficulty:      in.Difficulty,
		CookTimeMinutes: in.CookTimeMinutes,
		MainIngredients: strings.TrimSpace(in.MainIngredients),
		Taste:           strings.TrimSpace(in.Taste),
		Scene:           strings.TrimSpace(in.Scene),
		Note:            in.Note,
		Extra:           datatypes.JSON([]byte("{}")),
	}
	if in.IsFavorite != nil {
		dish.IsFavorite = *in.IsFavorite
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dish).Error; err != nil {
			return err
		}
		if err := s.syncTags(tx, dish.ID, dish, in.Tags); err != nil {
			return err
		}
		return s.bindCover(tx, dish, in.CoverAttachmentID)
	})
	if err != nil {
		return nil, err
	}
	return s.Get(dish.ID)
}

func (s *DishService) Update(id uint64, in model.DishCreateInput) (*model.Dish, error) {
	dish, err := s.dishes.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, fmt.Errorf("菜名为必填")
	}
	status := in.Status
	if status == "" {
		status = dish.Status
	}
	if !model.ValidDishStatus(status) {
		return nil, fmt.Errorf("无效的菜品状态")
	}
	dish.Name = name
	dish.CoverImageURL = strings.TrimSpace(in.CoverImageURL)
	dish.SourceURL = strings.TrimSpace(in.SourceURL)
	dish.SourcePlatform = strings.TrimSpace(in.SourcePlatform)
	dish.Description = strings.TrimSpace(in.Description)
	dish.Status = status
	dish.Rating = in.Rating
	dish.Difficulty = in.Difficulty
	dish.CookTimeMinutes = in.CookTimeMinutes
	dish.MainIngredients = strings.TrimSpace(in.MainIngredients)
	dish.Taste = strings.TrimSpace(in.Taste)
	dish.Scene = strings.TrimSpace(in.Scene)
	dish.Note = in.Note
	if in.IsFavorite != nil {
		dish.IsFavorite = *in.IsFavorite
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(dish).Error; err != nil {
			return err
		}
		if err := s.syncTags(tx, dish.ID, dish, in.Tags); err != nil {
			return err
		}
		return s.bindCover(tx, dish, in.CoverAttachmentID)
	})
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *DishService) Get(id uint64) (*model.Dish, error) {
	dish, err := s.dishes.GetByID(id)
	if err != nil {
		return nil, err
	}
	s.attachCover(dish)
	records, err := s.recs.ListByDish(id)
	if err != nil {
		return nil, err
	}
	if err := s.attachRecordPhotos(records); err != nil {
		return nil, err
	}
	dish.Records = records
	return dish, nil
}

func (s *DishService) List(q model.DishListQuery) (*model.PageResult[model.Dish], error) {
	items, total, err := s.dishes.List(q)
	if err != nil {
		return nil, err
	}
	s.attachCovers(items)
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	return &model.PageResult[model.Dish]{Items: items, Page: page, PageSize: size, Total: total}, nil
}

func (s *DishService) Delete(id uint64) error {
	_, err := s.dishes.GetByID(id)
	if err != nil {
		return err
	}
	return s.dishes.SoftDelete(id)
}

func (s *DishService) syncTags(tx *gorm.DB, dishID uint64, dish *model.Dish, extra []model.TagInput) error {
	var old []model.DishTag
	if err := tx.Where("dish_id = ?", dishID).Find(&old).Error; err != nil {
		return err
	}

	type key struct{ name, typ string }
	seen := map[key]*model.Tag{}
	add := func(name, typ string) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return nil
		}
		if typ == "" {
			typ = model.TagCustom
		}
		if !model.ValidTagType(typ) {
			typ = model.TagCustom
		}
		k := key{name, typ}
		if _, ok := seen[k]; ok {
			return nil
		}
		tag, err := s.tags.Upsert(tx, name, typ)
		if err != nil {
			return err
		}
		seen[k] = tag
		return nil
	}

	for _, n := range splitNames(dish.MainIngredients) {
		if err := add(n, model.TagIngredient); err != nil {
			return err
		}
	}
	for _, t := range extra {
		if t.Type == model.TagTaste || t.Type == model.TagScene {
			continue
		}
		if err := add(t.Name, t.Type); err != nil {
			return err
		}
	}

	tags := make([]model.Tag, 0, len(seen))
	ids := make([]uint64, 0, len(seen)+len(old))
	for _, t := range seen {
		tags = append(tags, *t)
		ids = append(ids, t.ID)
	}
	for _, o := range old {
		ids = append(ids, o.TagID)
	}
	if err := s.dishes.ReplaceTags(tx, dishID, tags); err != nil {
		return err
	}
	return s.tags.RecalcUsage(tx, ids)
}

func (s *DishService) bindCover(tx *gorm.DB, dish *model.Dish, attID *uint64) error {
	repo := repository.NewAttachmentRepo(tx)
	if attID != nil && *attID > 0 {
		if err := repo.Bind([]uint64{*attID}, model.BizDishCover, dish.ID); err != nil {
			return err
		}
	}
	if dish.CoverImageURL != "" {
		return repo.BindByURLs([]string{dish.CoverImageURL}, model.BizDishCover, dish.ID)
	}
	return nil
}

func (s *DishService) attachCover(dish *model.Dish) {
	atts, err := s.atts.ListByBiz(model.BizDishCover, dish.ID)
	if err != nil || len(atts) == 0 {
		return
	}
	if dish.CoverImageURL == "" {
		dish.CoverImageURL = atts[0].FileURL
	}
	if atts[0].ThumbnailURL != "" {
		dish.CoverThumbnailURL = atts[0].ThumbnailURL
	} else {
		dish.CoverThumbnailURL = atts[0].FileURL
	}
}

func (s *DishService) attachCovers(items []model.Dish) {
	if len(items) == 0 {
		return
	}
	ids := make([]uint64, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].ID)
	}
	atts, err := s.atts.ListByBizIDs(model.BizDishCover, ids)
	if err != nil {
		return
	}
	byDish := map[uint64]model.Attachment{}
	for _, a := range atts {
		if _, ok := byDish[a.BizID]; !ok {
			byDish[a.BizID] = a
		}
	}
	for i := range items {
		if a, ok := byDish[items[i].ID]; ok {
			if items[i].CoverImageURL == "" {
				items[i].CoverImageURL = a.FileURL
			}
			if a.ThumbnailURL != "" {
				items[i].CoverThumbnailURL = a.ThumbnailURL
			} else {
				items[i].CoverThumbnailURL = a.FileURL
			}
		}
	}
}

func (s *DishService) attachRecordPhotos(records []model.CookRecord) error {
	if len(records) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(records))
	for _, r := range records {
		ids = append(ids, r.ID)
	}
	atts, err := s.atts.ListByBizIDs(model.BizCookRecord, ids)
	if err != nil {
		return err
	}
	grouped := map[uint64][]model.Attachment{}
	for _, a := range atts {
		grouped[a.BizID] = append(grouped[a.BizID], a)
	}
	for i := range records {
		records[i].Photos = grouped[records[i].ID]
	}
	return nil
}

func (s *DishService) ThumbnailURL(dish *model.Dish) string {
	atts, err := s.atts.ListByBiz(model.BizDishCover, dish.ID)
	if err != nil || len(atts) == 0 {
		return dish.CoverImageURL
	}
	if atts[0].ThumbnailURL != "" {
		return atts[0].ThumbnailURL
	}
	return atts[0].FileURL
}
