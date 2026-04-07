package storage

import (
	"context"
	"fmt"
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPgCon(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := cfg.Dsn()
	poolCfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("error parse dsn %s %w", dsn, err)
	}

	poolCfg.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), poolCfg)

	if err != nil {
		return nil, fmt.Errorf("error connect pool %w", err)
	}

	return pool, nil
}
