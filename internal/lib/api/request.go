package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	httpReq   *http.Request
	validator *validator.Validate
}

func NewRequest(r *http.Request, v *validator.Validate) *Request {
	return &Request{
		httpReq:   r,
		validator: v,
	}
}

func (r *Request) Bind(dst any) error {
	if err := json.NewDecoder(r.httpReq.Body).Decode(dst); err != nil {
		return err
	}

	return r.validator.Struct(dst)
}
