package services

import (
	"context"
	"url-shortener/internal/models"
)

type clickWritter interface {
	Create(ctx context.Context, m *models.Click) error
}

type clickService struct {
	writer clickWritter
}

func NewClickService(w clickWritter) *clickService {
	return &clickService{
		writer: w,
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
