package handlers

import (
	"net/http"
	"url-shortener/internal/lib/api"
)

type stastHandler struct{}

func NewStatsHandler() *stastHandler {
	return &stastHandler{}
}

func (h *stastHandler) GetStatsByCode(r *http.Request) (api.Renderer, error) {
	return api.NewResponse(http.StatusOK, "handle", nil, nil)
}
