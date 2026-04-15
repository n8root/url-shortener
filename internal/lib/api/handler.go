package api

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"
)

type appError interface {
	GetStatus() int
	GetMessage() string
}

type errContainer interface {
	GetErrors() any
}

type action func(r *http.Request) (Renderer, error)

type Handler struct {
	Action action
}

func BindHandler(a action) http.HandlerFunc {
	h := Handler{Action: a}

	return h.Handle
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	res, err := h.Action(r)

	if err != nil {
		log.Printf("error: %v\n", err)
		log.Printf("Stack trace:\n%s\n", debug.Stack())

		var errRes = ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Server internal error",
		}

		var ae appError
		if ok := errors.As(err, &ae); ok {
			errRes.Status = ae.GetStatus()
			errRes.Message = ae.GetMessage()
		}

		var ec errContainer
		if ok := errors.As(err, &ec); ok {
			errRes.Errors = ec.GetErrors()
		}

		errRes.Render(w, r)
		return
	}

	res.Render(w, r)
}
