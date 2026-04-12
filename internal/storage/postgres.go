package storage

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgPool(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	const op = "storage.New"

	dsn := cfg.Dsn()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: error parse dsn: %w", op, err)
	}

	// Настройки пула
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	// Создаем пул
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: error create pool: %w", op, err)
	}

	// Проверяем, что база реально отвечает (Ping)
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: database is unreachable: %w", op, err)
	}

	return pool, nil
}
