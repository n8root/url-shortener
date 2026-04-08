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

	"github.com/go-playground/validator/v10"
)

func main() {
	//Чтобы запускать билд с аргументами --config=...
	cfgPath := flag.String("config", "configs/local.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)

	if err != nil {
		log.Fatalf("error get config %v", err)
	}

	serve(cfg)

}

func serve(cfg *config.Config) {
	v := validator.New()

	storage, err := storage.NewStorage(cfg)

	if err != nil {
		log.Fatalf("error init storage %v", err)
	}

	urlRepository := repositories.NewUrlRepository(storage)
	urlSerice := services.NewUrlService(urlRepository)
	urlHandler := handlers.NewUrlHandler(urlSerice)

	http.HandleFunc("POST /urls", api.Handle(v, urlHandler.Create))

	addr := cfg.Server.Addr()

	log.Printf("Server started on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
