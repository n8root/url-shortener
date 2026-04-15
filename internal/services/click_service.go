package services

import (
	"context"
	"url-shortener/internal/models"
)

type clickWriter interface {
	Create(ctx context.Context, m *models.Click) error
}

type clickReader interface {
	GetStatsByCode(ctx context.Context, code string) (*models.Stat, error)
}

type clickService struct {
	writer clickWriter
	reader clickReader
}

func NewClickService(w clickWriter, r clickReader) *clickService {
	return &clickService{
		writer: w,
		reader: r,
	}
}

func (s *clickService) Create(ctx context.Context, f *models.CreateClickForm) (*models.Click, error) {
	m := &models.Click{
		UrlID:     f.UrlID,
		IP:        f.IP,
		UserAgent: f.UserAgent,
		Refer:     f.Refer,
	}

	err := s.writer.Create(ctx, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *clickService) GetStatsByCode(ctx context.Context, code string) (*models.Stat, error) {
	m, err := s.reader.GetStatsByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return m, err
}
