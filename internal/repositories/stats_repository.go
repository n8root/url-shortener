package repositories

import (
	"context"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
)

type statsRepository struct {
	s *storage.Storage
}

func NewStatsRepository(s *storage.Storage) *statsRepository {
	return &statsRepository{s: s}
}

func (r *statsRepository) GetStatsByCode(ctx context.Context, code string) ([]models.Stat, error) {
	var col []models.Stat

	query := `
		SELECT 
			urls.code AS code,
			urls.original_url AS original_url,
			COUNT(*) as clicks_by_day,
			urls.created_at AS created_at,
			urls.created_at AS expires_at,
			clicks.created_at AS clicks_created_at
		FROM clicks
			LEFT JOIN urls ON urls.id=clicks.url_id
				WHERE urls.code=$1
					GROUP BY urls.created_at
	`

	rows, err := r.s.DB.Query(ctx, query, code)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clickMap := map[string][string]any{}

	for rows.Next() {
		var m models.Stat

		err := rows.Scan(
			&m.Code,
			&m.OriginalUrl,
			&m.ClicksByDay,
			&m.CreatedAt,
			&m.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}

		clickMap[m.Code]

		col = append(col, m)
	}

	return col, nil
}
