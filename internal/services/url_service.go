package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"url-shortener/internal/models"
)

type urlWritter interface {
	Save(context.Context, *models.Url) error
}

type urlReader interface {
	GetByCode(context.Context, string) (*models.Url, error)
	ExistsByCode(context.Context, string) (bool, error)
}

type UrlService struct {
	Writter urlWritter
	Reader  urlReader
}

func NewUrlService(writter urlWritter, reader urlReader) *UrlService {
	return &UrlService{
		Writter: writter,
		Reader:  reader,
	}
}

func (s *UrlService) Create(ctx context.Context, form models.CreateUrlForm) (*models.Url, error) {
	model := &models.Url{
		Code:        form.Alias,
		OriginalUrl: form.OriginalUrl,
		CustomAlias: form.Alias != "",
		ExpiresAt:   form.ExpiresAt,
	}

	if form.Alias == "" {
		model.Code = s.makeAlias(form.OriginalUrl)
	}

	exists, err := s.Reader.ExistsByCode(ctx, model.Code)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, models.EntityError{
			Status:  http.StatusConflict,
			Message: fmt.Sprintf("alias '%s' alredy exists", model.Code),
		}
	}

	err = s.Writter.Save(ctx, model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *UrlService) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	model, err := s.Reader.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	return model, err
}

func (s *UrlService) makeAlias(url string) string {
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	timestamp := time.Now().UnixNano()
	payload := fmt.Sprintf("%s%d%x", url, timestamp, randomBytes)

	hash := sha256.Sum256([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(hash[:])[:6]
}
