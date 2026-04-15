package validator

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	initCustomRules(v)

	return v
}

func initCustomRules(v *validator.Validate) {
	_ = v.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		if dateStr == "" {
			return true
		}

		layout := fl.Param()
		if layout == "" {
			layout = "2006-01-02"
		}

		_, err := time.Parse(layout, dateStr)
		return err == nil
	})

	_ = v.RegisterValidation("tomorrow", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		if dateStr == "" {
			return true
		}

		layout := fl.Param()
		if layout == "" {
			layout = "2006-01-02"
		}

		_, err := time.Parse(layout, dateStr)
		return err == nil
	})

	_ = v.RegisterValidation("excluded_aliases", func(fl validator.FieldLevel) bool {
		excludedAliases := []string{"admin", "api", "stats", "health", "urls"}
		str := fl.Field().String()

		if str == "" {
			return true
		}

		for _, a := range excludedAliases {
			if a == str {
				return false
			}
		}

		return true
	})

	_ = v.RegisterValidation("only_http", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()

		if str == "" {
			return true
		}

		if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
			return true
		}

		return false
	})
}
