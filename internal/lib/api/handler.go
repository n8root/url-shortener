package api

import (
	"context"
	"net/http"
	"url-shortener/internal/entities"

	"github.com/go-playground/validator/v10"
)

type HandleFunc[Req any, Res any] func(ctx context.Context, req Req) (Res, error)

func Handle[Req any, Res any](v *validator.Validate, logic HandleFunc[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodWithBody := r.Method == http.MethodPost ||
			r.Method == http.MethodPatch ||
			r.Method == http.MethodPut

		if methodWithBody && r.Body == http.NoBody || r.ContentLength == 0 {
			NewErrorResponse(&entities.EntityError{
				Status:  http.StatusBadRequest,
				Code:    "ERR_EMPTY_BODY",
				Message: "body required",
			}).Render(w)
			return
		}

		var body Req

		request := NewRequest(r, v)

		if err := request.Bind(&body); err != nil {
			NewErrorResponse(&entities.EntityError{
				Status:  http.StatusUnprocessableEntity,
				Code:    "ERR_VALIDATION",
				Message: "validation err",
			}).Render(w)
			return
		}

		res, err := logic(r.Context(), body)

		if err != nil {
			NewErrorResponse(err).Render(w)
			return
		}

		NewSuccessResponse(res).Render(w)
	}
}

func formatValidatonErrors(ve validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range ve {
		errs[f.Field()] = "failed on the '" + f.Tag() + "' tag"
	}

	return errs
}
