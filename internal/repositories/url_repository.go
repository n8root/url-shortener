package repositories

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/entities"
	"url-shortener/internal/storage"
)

type URLRepository interface {
	Save(ctx context.Context, entity *entities.Url) error
	GetByCode(ctx context.Context, code string) (*entities.Url, error)
}

type urlRepository struct {
	Storage *storage.Storage
}

func NewUrlRepository(storage *storage.Storage) *urlRepository {
	return &urlRepository{
		Storage: storage,
	}
}

func (r *urlRepository) Save(ctx context.Context, entity *entities.Url) error {
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

func (r *urlRepository) GetByCode(ctx context.Context, code string) (*entities.Url, error) {
	query := `
		SELECT id, code, original_url, custom_alias, created_at, expires_at, is_active
			FROM urls
				WHERE code = $1 AND is_active = true
	`

	var entity entities.Url

	err := r.Storage.DB.QueryRow(ctx, query, code).Scan(
		&entity.ID,
		&entity.Code,
		&entity.OriginalUrl,
		&entity.CustomAlias,
		&entity.CreatedAt,
		&entity.ExpiresAt,
		&entity.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.EntityError{
				Status:  404,
				Code:    "URL_NOT_FOUND",
				Message: "this alis is alredy taken",
			}
		}

		return nil, err
	}

	return &entity, nil
}
