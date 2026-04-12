package handlers

import (
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/models"
	"url-shortener/internal/services"

	"github.com/go-chi/chi"
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

func (h *UrlHandler) Create(r *http.Request) (api.Renderer, error) {
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

	res, err := api.NewResponse(http.StatusCreated, "Created", url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UrlHandler) RedirectByCode(r *http.Request) (api.Renderer, error) {
	code := chi.URLParam(r, "code")

	url, err := h.service.GetByCode(r.Context(), code)

	if err != nil {
		return nil, err
	}

	res, err := api.NewRedirectResponse(http.StatusTemporaryRedirect, url.OriginalUrl)
	if err != nil {
		return nil, err
	}

	return res, nil
}
