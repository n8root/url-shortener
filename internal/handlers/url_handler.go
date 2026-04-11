package handlers

import (
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/models"
	"url-shortener/internal/services"

	"github.com/go-playground/validator/v10"
)

type UrlHandler struct {
	service   *services.UrlService
	validator *validator.Validate
}

func NewUrlHandler(s *services.UrlService, v *validator.Validate) *UrlHandler {
	return &UrlHandler{
		service:   s,
		validator: v,
	}
}

func (h *UrlHandler) Create(r *http.Request) (*api.Response, error) {
	form := models.CreateUrlForm{}

	err := api.BindForm(r, &form)

	if err != nil {
		return nil, err
	}

	err = api.Validate(form, h.validator)

	if err != nil {
		return nil, err
	}

	url, err := h.service.Create(r.Context(), form)

	if err != nil {
		return nil, err
	}

	return &api.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    url,
	}, nil
}
