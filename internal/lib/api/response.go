package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
	Debug   any    `json:"debug,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
	Debug   any    `json:"debug,omitempty"`
}

type JSONResponse interface {
	GetStatus() int
}

func (r Response) GetStatus() int      { return r.Status }
func (r ErrorResponse) GetStatus() int { return r.Status }

func RenderJson(w http.ResponseWriter, r JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.GetStatus())
	_ = json.NewEncoder(w).Encode(r)
}
