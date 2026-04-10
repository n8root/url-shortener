package handlers

import (
	"context"
	"time"
	"url-shortener/internal/entities"
	"url-shortener/internal/services"

	"github.com/go-playground/validator/v10"
)

type CreateUrlRequest struct {
	Alias       string     `json:"alias" validate:"omitempty,alphanum,max=10"`
	OriginalUrl string     `json:"original_url" validate:"required,url"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type GetByCodeRequest struct {
	Code string
}

type UrlHandler struct {
	service   services.UrlSerice
	validator *validator.Validate
}

func NewUrlHandler(service services.UrlSerice) *UrlHandler {
	return &UrlHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *UrlHandler) Create(ctx context.Context, req CreateUrlRequest) (*entities.Url, error) {
	return h.service.Create(ctx, services.CreateParams(req))
}

func (h *UrlHandler) GetByCode(ctx context.Context, req GetByCodeRequest) (*entities.Url, error) {
	return h.service.GetByCode(ctx, req.Code)
}
