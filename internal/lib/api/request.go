package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"url-shortener/internal/models"

	"github.com/go-playground/validator/v10"
)

func Validate(s any, v *validator.Validate) error {
	if err := v.Struct(s); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return models.EntityError{
				Status:  http.StatusUnprocessableEntity,
				Message: "Validation error",
				Errors:  formatValidatonErrors(validateErrs),
			}
		}

		return err
	}

	return nil
}

func BindForm(r *http.Request, dst any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return models.EntityError{
				Status:  http.StatusBadRequest,
				Message: "Empty payload",
			}
		}

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return models.EntityError{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Malformed JSON at position %d", syntaxErr.Offset),
			}
		}

		return err
	}

	return nil
}

func formatValidatonErrors(v validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range v {
		errs[f.Field()] = "failed on the '" + f.Tag() + "' tag"
	}

	return errs
}
