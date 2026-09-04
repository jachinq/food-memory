package service

import (
	"errors"

	"food-memory/internal/model"
	"food-memory/internal/repository"

	"gorm.io/gorm"
)

type HomeService struct {
	home *repository.HomeRepo
	atts *repository.AttachmentRepo
}

func NewHomeService(home *repository.HomeRepo, atts *repository.AttachmentRepo) *HomeService {
	return &HomeService{home: home, atts: atts}
}

func (s *HomeService) Summary() (*model.HomeSummary, error) {
	stats, err := s.home.Stats()
	if err != nil {
		return nil, err
	}
	recent, err := s.home.Recent(6)
	if err != nil {
		return nil, err
	}
	overdue, err := s.home.OverdueHighRating(6)
	if err != nil {
		return nil, err
	}
	want, err := s.home.WantToCook(6)
	if err != nil {
		return nil, err
	}
	random, err := s.home.RandomOld()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	fillThumbnails := func(items []model.Dish) {
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
		byID := map[uint64]model.Attachment{}
		for _, a := range atts {
			if _, ok := byID[a.BizID]; !ok {
				byID[a.BizID] = a
			}
		}
		for i := range items {
			if a, ok := byID[items[i].ID]; ok {
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
	fillThumbnails(recent)
	fillThumbnails(overdue)
	fillThumbnails(want)
	if random != nil {
		atts, _ := s.atts.ListByBiz(model.BizDishCover, random.ID)
		if len(atts) > 0 {
			if random.CoverImageURL == "" {
				random.CoverImageURL = atts[0].FileURL
			}
			if atts[0].ThumbnailURL != "" {
				random.CoverThumbnailURL = atts[0].ThumbnailURL
			}
		}
	}
	return &model.HomeSummary{
		Stats:             stats,
		Recent:            recent,
		OverdueHighRating: overdue,
		RandomOld:         random,
		WantToCook:        want,
	}, nil
}
