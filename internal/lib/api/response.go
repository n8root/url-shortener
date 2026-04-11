package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	nurl "net/url"
)

type Renderer interface {
	Render(w http.ResponseWriter)
}

type statuser interface {
	GetStatus() int
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
}

func NewResponse(status int, message string, data any, meta any) (*Response, error) {
	if status < 100 || status > 511 {
		return nil, fmt.Errorf("status %d is incorrect", status)
	}

	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}, nil
}

func (r Response) GetStatus() int {
	return r.Status
}

func (r *Response) Render(w http.ResponseWriter) {
	renderJson(w, r)
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (r ErrorResponse) GetStatus() int {
	return r.Status
}

func NewErrorResponse(status int, message string, errors any) (*ErrorResponse, error) {
	if status < 400 || status > 511 {
		return nil, fmt.Errorf("status %d is incorrect", status)
	}

	return &ErrorResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	}, nil
}

func (r *ErrorResponse) Render(w http.ResponseWriter) {
	renderJson(w, r)
}

type RedirectResponse struct {
	Status int
	Url    string
}

func NewRedirectResponse(status int, url string) (*RedirectResponse, error) {
	if status < 300 || status > 308 {
		return nil, errors.New("not allowed redirect status")
	}

	if _, err := nurl.ParseRequestURI(url); err != nil {
		return nil, err
	}

	return &RedirectResponse{
		Status: status,
		Url:    url,
	}, nil
}

func (r *RedirectResponse) Render(w http.ResponseWriter) {
	http.Redirect(w, nil, r.Url, r.Status)
}

func renderJson(w http.ResponseWriter, r statuser) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.GetStatus())
	_ = json.NewEncoder(w).Encode(r)
}
