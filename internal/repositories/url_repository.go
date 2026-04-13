package repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5"
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

	err := r.Storage.DB.QueryRow(ctx, query, code).Scan(
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
	}

	return nil
}
