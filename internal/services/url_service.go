package services

import (
	"context"
	"errors"
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

type urlDeleter interface {
	DeleteByCode(context.Context, string) error
}

type urlService struct {
	writter urlWritter
	reader  urlReader
	deleter urlDeleter
}

func NewUrlService(w urlWritter, r urlReader, d urlDeleter) *urlService {
	return &urlService{writter: w, reader: r, deleter: d}
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
		err := model.MakeCode()
		if err != nil {
			return nil, err
		}
	}

	exists, err := s.reader.ExistsByCode(ctx, model.Code)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, models.EntityError{
			Status:  http.StatusConflict,
			Message: fmt.Sprintf("alias '%s' alredy exists", model.Code),
		}
	}

	err = s.writter.Save(ctx, model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *urlService) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	model, err := s.reader.GetByCode(ctx, code)

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

func (s *urlService) DeleteByCode(ctx context.Context, code string) error {
	if err := s.deleter.DeleteByCode(ctx, code); err != nil {
		return err
	}

	return nil
}

func (s *urlService) generateCode(ctx context.Context, url *models.Url) error {
	iteration := 3

	for iteration != 0 {
		err := url.MakeCode()

		if err != nil {
			return err
		}

		exists, err := s.reader.ExistsByCode(ctx, url.Code)
		if err != nil {
			return err
		}

		if !exists {
			return nil
		}

		iteration--
	}

	return errors.New("generate alias failed")
}
