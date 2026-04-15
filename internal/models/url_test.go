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
			name: "success alias",
			input: CreateUrlForm{
				Alias:       "test",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: false,
		},
		{
			name: "fail alias with -",
			input: CreateUrlForm{
				Alias:       "test-",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with 2 len",
			input: CreateUrlForm{
				Alias:       "te",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with 33 len",
			input: CreateUrlForm{
				Alias:       "te",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with spaces",
			input: CreateUrlForm{
				Alias:       "te  f",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with cyrillic",
			input: CreateUrlForm{
				Alias:       "Якот",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with services words",
			input: CreateUrlForm{
				Alias:       "health",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with not alphanum",
			input: CreateUrlForm{
				Alias:       "1221$3",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail alias with not alphanum",
			input: CreateUrlForm{
				Alias:       "1221$3",
				OriginalUrl: "https://test.com",
				ExpiresAt:   "2026-01-01",
			},
			wantErr: true,
		},
		{
			name: "fail expires at",
			input: CreateUrlForm{
				Alias:       "1221$3",
				OriginalUrl: "https://test.com",
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

func provideCodes(size int) map[string]string {
	codes := make(map[string]string, size)

	for i := 0; i < size; i++ {
		m := Url{OriginalUrl: "https://test.com"}
		_ = m.MakeCode()
		codes[m.Code] = m.Code
	}

	return codes
}

func TestMakeUniqCode(t *testing.T) {
	t.Run("All 1000 codes should be uniq", func(t *testing.T) {
		codes := provideCodes(1000)
		assert.Len(t, codes, 1000)
	})
}

func TestMakeCodeWithCorrectLength(t *testing.T) {
	t.Run("All 1000 codes should be 6 >= code <=", func(t *testing.T) {
		codes := provideCodes(1000)
		for _, code := range codes {
			assert.True(t, 6 <= len(code) && 8 >= len(code))
		}
	})
}
