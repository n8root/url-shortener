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

type urlService struct {
	Writter urlWritter
	Reader  urlReader
}

func NewUrlService(writter urlWritter, reader urlReader) *urlService {
	return &urlService{
		Writter: writter,
		Reader:  reader,
	}
}

func (s *urlService) Create(ctx context.Context, f *models.CreateUrlForm) (*models.Url, error) {
	t, _ := time.Parse(time.DateOnly, f.ExpiresAt)

	model := &models.Url{
		Code:        f.Alias,
		OriginalUrl: f.OriginalUrl,
		CustomAlias: f.Alias != "",
		ExpiresAt:   &t,
	}

	if f.Alias == "" {
		model.Code = s.makeAlias(f.OriginalUrl)
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

func (s *urlService) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	model, err := s.Reader.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	if model.IsExpired() {
		return nil, models.EntityError{
			Status:  http.StatusGone,
			Message: fmt.Sprintf("alias '%s' is expired", model.Code),
		}
	}

	return model, err
}

func (s *urlService) makeAlias(url string) string {
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	timestamp := time.Now().UnixNano()
	payload := fmt.Sprintf("%s%d%x", url, timestamp, randomBytes)

	hash := sha256.Sum256([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(hash[:])[:6]
}
