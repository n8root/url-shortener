package main

import (
	"flag"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/validator"

	"github.com/go-chi/chi"
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
	validator := validator.NewValidator()

	router := chi.NewRouter()

	storage, err := storage.NewStorage(cfg)

	if err != nil {
		log.Fatalf("failed initialization storage %v", err)
	}

	clickRepo := repositories.NewClickRepository(storage)
	clickSvc := services.NewClickService(clickRepo, clickRepo)

	urlRepo := repositories.NewUrlRepository(storage)
	urlSvc := services.NewUrlService(urlRepo, urlRepo, urlRepo)

	urlHandler := handlers.NewUrlHandler(urlSvc, clickSvc, validator)
	statsHandler := handlers.NewStatsHandler(clickSvc)

	router.Route("/urls", func(r chi.Router) {
		r.Post("/", api.BindHandler(urlHandler.Create))
		r.Get("/{code}", api.BindHandler(urlHandler.RedirectByCode))
		r.Delete("/{code}", api.BindHandler(urlHandler.DeleteByCode))
	})

	router.Route("/stats", func(r chi.Router) {
		r.Get("/{code}", api.BindHandler(statsHandler.GetStatsByCode))
	})

	addr := cfg.Server.Addr()

	log.Printf("Server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
