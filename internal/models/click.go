package models

import (
	"time"
)

type Click struct {
	ID        int        `json:"id"`
	UrlID     int        `json:"url_id"`
	IP        string     `json:"ip"`
	UserAgent string     `json:"user_agent"`
	Refer     string     `json:"refer"`
	CreatedAt *time.Time `json:"created_at"`
}

type CreateClickForm struct {
	UrlID     int    `json:"url_id" validate:"required,alphanum"`
	IP        string `json:"ip" validate:"required,ip"`
	UserAgent string `json:"user_agent" validate:"required,string"`
	Refer     string `json:"refer" validate:"required,url"`
}
