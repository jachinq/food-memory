package service

import (
	"fmt"
	"strings"

	"food-memory/internal/model"
	"food-memory/internal/repository"
)

type TagService struct {
	tags *repository.TagRepo
}

func NewTagService(tags *repository.TagRepo) *TagService {
	return &TagService{tags: tags}
}

func (s *TagService) List(tagType, keyword string) ([]model.Tag, error) {
	return s.tags.List(tagType, keyword)
}

func (s *TagService) Create(name, tagType string) (*model.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("标签名不能为空")
	}
	if tagType == "" {
		tagType = model.TagCustom
	}
	if !model.ValidTagType(tagType) {
		return nil, fmt.Errorf("无效的标签类型")
	}
	tag := &model.Tag{Name: name, Type: tagType}
	if err := s.tags.Create(tag); err != nil {
		existing, err2 := s.tags.Upsert(nil, name, tagType)
		if err2 != nil {
			return nil, err
		}
		return existing, nil
	}
	return tag, nil
}
