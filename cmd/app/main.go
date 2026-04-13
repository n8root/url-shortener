package main

import (
	"flag"
	"log"
	"net/http"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

func main() {
	//Чтобы запускать билд с аргументами --config=...
	cfgPath := flag.String("config", "configs/local.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)

	if err != nil {
		log.Fatalf("failed initialization config %v", err)
	}

	serve(cfg)

}

func serve(cfg *config.Config) {
	validator := validator.New()

	registerCustomRules(validator)

	router := chi.NewRouter()

	storage, err := storage.NewStorage(cfg)

	if err != nil {
		log.Fatalf("failed initialization storage %v", err)
	}

	clickRepo := repositories.NewClickRepository(storage)
	clickSvc := services.NewClickService(clickRepo)

	urlRepo := repositories.NewUrlRepository(storage)
	urlSvc := services.NewUrlService(urlRepo, urlRepo)

	urlHandler := handlers.NewUrlHandler(urlSvc, clickSvc, validator)

	router.Route("/urls", func(r chi.Router) {
		r.Post("/", api.BindHandler(urlHandler.Create))
		r.Get("/{code}", api.BindHandler(urlHandler.RedirectByCode))
	})

	addr := cfg.Server.Addr()

	log.Printf("Server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func registerCustomRules(v *validator.Validate) {
	_ = v.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		if dateStr == "" {
			return true
		}

		layout := fl.Param()
		if layout == "" {
			layout = "2006-01-02"
		}

		_, err := time.Parse(layout, dateStr)
		return err == nil
	})

	_ = v.RegisterValidation("tomorrow", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		if dateStr == "" {
			return true
		}

		layout := fl.Param()
		if layout == "" {
			layout = "2006-01-02"
		}

		_, err := time.Parse(layout, dateStr)
		return err == nil
	})
}
