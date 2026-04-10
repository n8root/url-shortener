package api

import (
	"errors"
	"fmt"
	"net/http"
)

type appError interface {
	GetStatus() int
	GetCode() string
	GetMessage() string
}

type errContainer interface {
	GetErrors() any
}

type handler func(r *http.Request) (*Response, error)

func Handle(h handler, w http.ResponseWriter, r *http.Request) {
	res, err := h(r)

	if err != nil {
		var errRes = ErrorResponse{
			Code:    "SERVER_INTERNAL_ERR",
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Server internal error %v", err),
		}

		var ae appError
		if ok := errors.As(err, &ae); ok {
			errRes.Code = ae.GetCode()
			errRes.Status = ae.GetStatus()
			errRes.Message = ae.GetMessage()
		}

		var ec errContainer
		if ok := errors.As(err, &ec); ok {
			errRes.Errors = ec.GetErrors()
		}

		RenderJson(w, errRes)
		return
	}

	RenderJson(w, res)
}
