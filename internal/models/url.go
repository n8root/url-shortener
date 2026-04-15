package models

import (
	"crypto/rand"
	"errors"
	"time"
)

type Url struct {
	ID          int        `json:"id" db:"id"`
	Code        string     `json:"code" db:"code"`
	OriginalUrl string     `json:"original_url" db:"original_id"`
	CustomAlias bool       `json:"custom_alias" db:"custom_alias"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
}

func (m Url) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}

	return time.Now().UTC().After(m.ExpiresAt.UTC())
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (m *Url) MakeCode() error {
	if m.OriginalUrl == "" {
		return errors.New("empty original url")
	}

	// 1. Генерируем случайное смещение для длины (0, 1 или 2)
	lenBuf := make([]byte, 1)
	if _, err := rand.Read(lenBuf); err != nil {
		return err
	}
	// Длина будет 6 + (0..2) = 6, 7 или 8
	length := 6 + int(lenBuf[0]%3)

	// 2. Генерируем случайные байты для выбора символов
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return err
	}

	// 3. Собираем итоговый код
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphabet[randomBytes[i]%byte(len(alphabet))]
	}

	m.Code = string(result)
	return nil
}

type CreateUrlForm struct {
	Alias       string `json:"alias" validate:"omitempty,alphanum,max=32,min=2,excluded_aliases"`
	OriginalUrl string `json:"original_url" validate:"required,url,only_http,max=2048"`
	ExpiresAt   string `json:"expires_at" validate:"omitempty,date_format=2006-01-02"`
}
