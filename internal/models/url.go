package models

import "time"

type Url struct {
	ID          int        `json:"id" db:"id"`
	Code        string     `json:"code" db:"code"`
	OriginalUrl string     `json:"original_url" db:"original_id"`
	CustomAlias bool       `json:"custom_alias" db:"custom_alias"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
}

type CreateUrlForm struct {
	Alias       string     `json:"alias" validate:"omitempty,alphanum,max=10"`
	OriginalUrl string     `json:"original_url" validate:"required,url"`
	ExpiresAt   *time.Time `json:"expires_at" validate:"omitempty,date_format=2006-01-02"`
}
