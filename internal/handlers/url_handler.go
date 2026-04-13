package handlers

import (
	"context"
	"log"
	"net"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type urlService interface {
	Create(ctx context.Context, f *models.CreateUrlForm) (*models.Url, error)
	GetByCode(ctx context.Context, code string) (*models.Url, error)
}

type clickService interface {
	Create(ctx context.Context, f *models.CreateClickForm) (*models.Click, error)
}

type urlHandler struct {
	urlServicve  urlService
	clickService clickService
	validator    *validator.Validate
}

func NewUrlHandler(
	urlSerice urlService,
	clickService clickService,
	v *validator.Validate,
) *urlHandler {
	return &urlHandler{
		urlServicve:  urlSerice,
		clickService: clickService,
		validator:    v,
	}
}

func (h *urlHandler) Create(r *http.Request) (api.Renderer, error) {
	f := models.CreateUrlForm{}

	err := api.BindForm(r, &f)
	if err != nil {
		return nil, err
	}

	err = api.Validate(f, h.validator)
	if err != nil {
		return nil, err
	}

	url, err := h.urlServicve.Create(r.Context(), &f)
	if err != nil {
		return nil, err
	}

	res, err := api.NewResponse(http.StatusCreated, "Created", url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *urlHandler) RedirectByCode(r *http.Request) (api.Renderer, error) {
	code := chi.URLParam(r, "code")

	url, err := h.urlServicve.GetByCode(r.Context(), code)
	if err != nil {
		return nil, err
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	refer := r.Referer()
	userAgent := r.UserAgent()

	h.createClickAsync(
		url.ID,
		url.Code,
		ip,
		refer,
		userAgent,
	)

	return api.NewRedirectResponse(http.StatusTemporaryRedirect, url.OriginalUrl)
}

func (h *urlHandler) createClickAsync(urlID int, code, ip, refer, ua string) {
	go func() {
		ctx := context.Background()

		f := &models.CreateClickForm{
			UrlID:     urlID,
			IP:        ip,
			UserAgent: ua,
			Refer:     refer,
		}

		if _, err := h.clickService.Create(ctx, f); err != nil {
			log.Printf("ERROR: to register click for code %s: %v", code, err)
		}
	}()
}
