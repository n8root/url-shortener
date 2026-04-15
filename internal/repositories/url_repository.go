package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type urlRepository struct {
	Storage *storage.Storage
}

func NewUrlRepository(storage *storage.Storage) *urlRepository {
	return &urlRepository{
		Storage: storage,
	}
}

func (r *urlRepository) Save(ctx context.Context, entity *models.Url) error {
	query := `
		INSERT INTO urls (code, original_url, custom_alias, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, is_active
	`

	err := r.Storage.DB.QueryRow(
		ctx,
		query,
		entity.Code,
		entity.OriginalUrl,
		entity.CustomAlias,
		entity.ExpiresAt,
	).Scan(
		&entity.ID,
		&entity.CreatedAt,
		&entity.IsActive,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *urlRepository) GetByCode(ctx context.Context, code string) (*models.Url, error) {
	var model models.Url

	cached, err := r.Storage.Cache.Get(ctx, "links:"+code).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err != redis.Nil {
		if err := json.Unmarshal([]byte(cached), &model); err != nil {
			return nil, err
		}

		return &model, nil
	}

	query := `
		SELECT id, code, original_url, custom_alias, created_at, expires_at, is_active
			FROM urls
				WHERE code = $1 AND is_active = true
	`

	err = r.Storage.DB.QueryRow(ctx, query, code).Scan(
		&model.ID,
		&model.Code,
		&model.OriginalUrl,
		&model.CustomAlias,
		&model.CreatedAt,
		&model.ExpiresAt,
		&model.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.EntityError{
				Status:  404,
				Message: fmt.Sprintf("url by code %s not found", code),
			}
		}

		return nil, err
	}

	ttl := time.Until(*model.ExpiresAt)

	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	stat := r.Storage.Cache.Set(ctx, "links:"+model.Code, data, ttl)

	if stat.Err() != nil {
		return nil, stat.Err()
	}

	return &model, nil
}

func (r *urlRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	exists := false

	query := `SELECT EXISTS (SELECT 1 FROM urls WHERE code=$1)`

	if err := r.Storage.DB.QueryRow(ctx, query, code).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *urlRepository) DeleteByCode(ctx context.Context, code string) error {
	query := `DELETE FROM urls WHERE code=$1`

	res, err := r.Storage.DB.Exec(ctx, query, code)
	if err != nil {
		return err
	}

	if res.RowsAffected() < 1 {
		return models.EntityError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("url by '%s' not found", code),
		}
	} else {
		r.Storage.Cache.Del(ctx, "links:"+code)
	}

	return nil
}
