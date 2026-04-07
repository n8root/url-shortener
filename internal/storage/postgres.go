package storage

import (
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConnPool(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {

}
