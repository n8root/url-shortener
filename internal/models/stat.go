package models

import "time"

type Stat struct {
	Code        string         `json:"code"`
	OriginalUrl string         `json:"original_url"`
	TotalClicks int            `json:"total_clicks"`
	ClicksByDay map[string]any `json:"clicks_by_day"`
	CreatedAt   time.Time      `json:"created_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`
}
