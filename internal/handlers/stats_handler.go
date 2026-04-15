package handlers

import (
	"context"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/models"

	"github.com/go-chi/chi"
)

type statsService interface {
	GetStatsByCode(ctx context.Context, code string) (*models.Stat, error)
}

type stastHandler struct {
	service statsService
}

func NewStatsHandler(s statsService) *stastHandler {
	return &stastHandler{
		service: s,
	}
}

func (h *stastHandler) GetStatsByCode(r *http.Request) (api.Renderer, error) {
	code := chi.URLParam(r, "code")

	m, err := h.service.GetStatsByCode(r.Context(), code)

	if err != nil {
		return nil, err
	}

	return api.NewResponse(http.StatusOK, "Success", m, nil)
}
