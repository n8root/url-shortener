package models

import (
	"strings"
	"testing"
	"url-shortener/internal/validator"

	"github.com/stretchr/testify/assert"
)

func TestCreateUrlFormValidation(t *testing.T) {
	v := validator.NewValidator()

	tests := []struct {
		name    string
		input   CreateUrlForm
		wantErr bool
	}{
		{
			name: "pass original_url with http",
			input: CreateUrlForm{
				Alias:       "123xbc",
				OriginalUrl: "http://test.url",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: false,
		},
		{
			name: "pass original_url with https",
			input: CreateUrlForm{
				Alias:       "123xbc",
				OriginalUrl: "https://test.url",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: false,
		},
		{
			name: "fail original_url with ftp",
			input: CreateUrlForm{
				Alias:       "123xbc",
				OriginalUrl: "ftp://test.url",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail with empty original_url",
			input: CreateUrlForm{
				Alias:       "123xbc",
				OriginalUrl: "",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail with max=2048 original_url",
			input: CreateUrlForm{
				Alias:       "123xbc",
				OriginalUrl: "https://test." + strings.Repeat("c", 2048-12),
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with not alphanum",
			input: CreateUrlForm{
				Alias:       "1221$3",
				OriginalUrl: "test.url",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail expires at",
			input: CreateUrlForm{
				Alias:       "1221$3",
				OriginalUrl: "test.url",
				ExpiresAt:   "123k,ofewf",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMakeCode(t *testing.T) {

}
