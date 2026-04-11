package repositories

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
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
	query := `
		SELECT id, code, original_url, custom_alias, created_at, expires_at, is_active
			FROM urls
				WHERE code = $1 AND is_active = true
	`

	var model models.Url

	err := r.Storage.DB.QueryRow(ctx, query, code).Scan(&model)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.EntityError{
				Status:  404,
				Message: "this alis is alredy taken",
			}
		}

		return nil, err
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
