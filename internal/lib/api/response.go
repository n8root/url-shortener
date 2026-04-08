package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
	Data    any    `json:"data,omitempty"`
	Debug   any    `json:"debug,omitempty"`
}

func NewErrorResponse(err error) *Response {
	res := &Response{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		//Debug:   err.Error(),
	}

	type coder interface{ GetCode() string }
	type statuser interface{ GetStatus() int }
	type messeger interface{ GetMessage() string }
	type errorer interface{ GetErrors() any }

	if c, ok := err.(coder); ok {
		res.Code = c.GetCode()
	}

	if s, ok := err.(statuser); ok {
		res.Status = s.GetStatus()
	}

	if m, ok := err.(messeger); ok {
		res.Message = m.GetMessage()
	}

	if e, ok := err.(errorer); ok {
		res.Errors = e.GetErrors()
	}

	return res
}

func NewSuccessResponse(data any) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

func (r *Response) Render(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	_ = json.NewEncoder(w).Encode(r)
}
