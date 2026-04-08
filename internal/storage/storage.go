package storage

import (
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	DB    *pgxpool.Pool
	Cache *redis.Client
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	pg, err := NewPgCon(&cfg.Database)
	if err != nil {
		return nil, err
	}

	rd := NewRedisCon(&cfg.Cache)

	storage := &Storage{
		DB:    pg,
		Cache: rd,
	}

	return storage, nil
}
