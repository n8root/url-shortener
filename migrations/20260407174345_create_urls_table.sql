-- +goose Up

CREATE TABLE urls(
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    original_url TEXT NOT NULL,
    custom_alias BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    expires_at TIMESTAMP
    is_active BOOLEAN NOT NULL DEFAULT TRUE
)

CREATE INDEX urls_code_idx ON urls(code)
CREATE INDEX urls_expires_at_idx ON urls(expires_at)
CREATE INDEX urls_created_at_idx ON urls(created_at)
