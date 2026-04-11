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
	router := chi.NewRouter()

	storage, err := storage.NewStorage(cfg)

	if err != nil {
		log.Fatalf("failed initialization storage %v", err)
	}

	urlRepo := repositories.NewUrlRepository(storage)
	urlSerice := services.NewUrlService(urlRepo, urlRepo)
	urlHandler := handlers.NewUrlHandler(urlSerice, validator)

	router.Route("/urls", func(r chi.Router) {
		r.Post("/", api.BindHandler(urlHandler.Create))
	})

	addr := cfg.Server.Addr()

	log.Printf("Server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
