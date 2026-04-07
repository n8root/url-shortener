-- +goose Up

CREATE TABLE clicks(
    id SERIAL PRIMARY KEY,
    url_id BIGINT NOT NULL REFERENCES urls(id) ON DELEATE CASCADE,
    ip INET NOT NULL
    user_agent TEXT
    refer TEXT
    created_at TIMESTAMP NOT NULL DEFAULT now()
)

CREATE INDEX clicks_url_id_create_at_idx ON urls(url_id, created_at)
