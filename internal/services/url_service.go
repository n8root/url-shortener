package services

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/entities"
	"url-shortener/internal/repositories"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("url not found")
	ErrAliasAlredyExists = errors.New("alias alredy exists")
)

type UrlSerice interface {
	Create(ctx context.Context, params CreateParams) (*entities.Url, error)
	GetByCode(ctx context.Context, code string) (*entities.Url, error)
}

type urlService struct {
	repo repositories.URLRepository
}

type CreateParams struct {
	Alias       string
	OriginalUrl string
	ExpiresAt   *time.Time
}

func NewUrlService(repo repositories.URLRepository) UrlSerice {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) Create(ctx context.Context, params CreateParams) (*entities.Url, error) {
	entity := &entities.Url{
		Code:        params.Alias,
		OriginalUrl: params.OriginalUrl,
		CustomAlias: params.Alias != "",
		ExpiresAt:   params.ExpiresAt,
	}

	if params.Alias == "" {
		entity.Code = uuid.New().String()
	}

	err := s.repo.Save(ctx, entity)

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *urlService) GetByCode(ctx context.Context, code string) (*entities.Url, error) {
	entity, err := s.repo.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	return entity, err
}
