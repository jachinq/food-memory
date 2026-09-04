package service

import (
	"mime/multipart"

	"food-memory/internal/config"
	"food-memory/internal/model"
	"food-memory/internal/repository"
	"food-memory/internal/storage"
)

type UploadService struct {
	cfg  *config.Config
	fs   *storage.Local
	atts *repository.AttachmentRepo
}

func NewUploadService(cfg *config.Config, fs *storage.Local, atts *repository.AttachmentRepo) *UploadService {
	return &UploadService{cfg: cfg, fs: fs, atts: atts}
}

func (s *UploadService) Save(file multipart.File, header *multipart.FileHeader, bizType string, bizID uint64) (*model.Attachment, error) {
	if bizType == "" {
		bizType = model.BizDishCover
	}
	saved, err := s.fs.Save(file, header, s.cfg.MaxUploadBytes)
	if err != nil {
		return nil, err
	}
	att := &model.Attachment{
		BizType:      bizType,
		BizID:        bizID,
		FileName:     saved.FileName,
		FileURL:      saved.FileURL,
		ThumbnailURL: saved.ThumbnailURL,
		MimeType:     saved.MimeType,
		FileSize:     saved.FileSize,
		Width:        saved.Width,
		Height:       saved.Height,
	}
	if err := s.atts.Create(att); err != nil {
		return nil, err
	}
	return att, nil
}
