package models

import "time"

type Url struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	OriginalUrl string     `json:"original_url"`
	CustomAlias bool       `json:"custom_alias"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"` // Используем указатель, чтобы в JSON был null
	IsActive    bool       `json:"is_active"`
}

type CreateUrlForm struct {
	Alias       string     `json:"alias" validate:"omitempty,alphanum,max=10"`
	OriginalUrl string     `json:"original_url" validate:"required,url"`
	ExpiresAt   *time.Time `json:"expires_at" validate:"omitempty,gt"`
}
