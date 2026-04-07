package storage

import (
	"url-shortener/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisCon(cfg *config.CacheConfig) *redis.Client {
	options := &redis.Options{
		DB:   cfg.DB,
		Addr: cfg.Addr(),
	}

	return redis.NewClient(options)
}
