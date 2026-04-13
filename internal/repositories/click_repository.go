package repositories

import (
	"context"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
)

type clickRepository struct {
	s *storage.Storage
}

func NewClickRepository(s *storage.Storage) *clickRepository {
	return &clickRepository{s: s}
}

func (r *clickRepository) Create(ctx context.Context, model *models.Click) error {
	query := `
		INSERT INTO clicks (url_id, ip, user_agent, refer)
			VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.s.DB.QueryRow(
		ctx,
		query,
		model.UrlID,
		model.IP,
		model.UserAgent,
		model.Refer,
	).Scan(
		&model.ID,
		&model.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
