package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/validator"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	validator := validator.NewValidator()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	storage, err := storage.NewStorage(cfg, ctx)

	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища %v", err)
	}

	defer storage.DB.Close()

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

	srv := &http.Server{Addr: addr, Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	log.Printf("Сервер запущен на %s", addr)

	<-ctx.Done()
	log.Println("Получен сигнал остановки, завершаем работу...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
