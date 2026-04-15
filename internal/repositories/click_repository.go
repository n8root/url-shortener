package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

func (r *clickRepository) GetStatsByCode(ctx context.Context, code string) (*models.Stat, error) {
	query := `
		SELECT 
			u.code,
			u.original_url,
			u.created_at,
			u.expires_at,
			DATE(c.created_at) as click_date,
			COUNT(c.id) as daily_count
		FROM urls u
			LEFT JOIN clicks c ON u.id = c.url_id
				WHERE u.code = $1
					GROUP BY u.id, click_date
						ORDER BY click_date DESC;
	`

	rows, err := r.s.DB.Query(ctx, query, code)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stat *models.Stat

	for rows.Next() {
		var clickDate *time.Time
		var dailyCount int

		if stat == nil {
			stat = &models.Stat{ClicksByDay: []models.DailyStat{}}
			err = rows.Scan(
				&stat.Code,
				&stat.OriginalUrl,
				&stat.CreatedAt,
				&stat.ExpiresAt,
				&clickDate,
				&dailyCount,
			)
		} else {
			err = rows.Scan(nil, nil, nil, nil, &clickDate, &dailyCount)
		}

		if err != nil {
			return nil, err
		}

		if clickDate != nil {
			stat.TotalClicks += dailyCount
			stat.ClicksByDay = append(stat.ClicksByDay, models.DailyStat{
				Date:  clickDate.Format("2006-01-02"),
				Count: dailyCount,
			})
		}
	}

	if stat == nil {
		return nil, models.EntityError{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("stats by code '%s' not found", code),
		}
	}

	return stat, nil
}
