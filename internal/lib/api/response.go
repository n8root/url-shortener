package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type entityError interface {
	GetStatus() int
	GetMessage() string
	GetCode() string
}

type Response struct {
	Status  int    `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewErrorResponse(err error) *Response {
	var ve validator.ValidationErrors

	if ae, ok := err.(entityError); ok {
		return &Response{
			Status:  ae.GetStatus(),
			Code:    ae.GetCode(),
			Message: ae.GetMessage(),
		}
	}

	if errors.As(err, &ve) {
		return &Response{
			Status:  http.StatusInternalServerError,
			Message: "validation error",
			Data:    formatValidatonErrors(ve),
		}
	}

	return &Response{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
	}
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

func formatValidatonErrors(ve validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range ve {
		errs[f.Field()] = "failed on the '" + f.Tag() + "' tag"
	}

	return errs
}
