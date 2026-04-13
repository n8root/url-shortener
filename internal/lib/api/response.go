package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	nurl "net/url"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request)
}

type statuser interface {
	GetStatus() int
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
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

func (res Response) GetStatus() int {
	return res.Status
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) {
	renderJson(w, res)
}

type NoContentReponce struct{}

func (res NoContentReponce) GetStatus() int {
	return http.StatusNoContent
}

func (res NoContentReponce) Render(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(res.GetStatus())
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (res ErrorResponse) GetStatus() int {
	return res.Status
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

func (res *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) {
	renderJson(w, res)
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

func (res *RedirectResponse) Render(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, res.Url, res.Status)
}

func renderJson(w http.ResponseWriter, r statuser) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.GetStatus())
	_ = json.NewEncoder(w).Encode(r)
}
