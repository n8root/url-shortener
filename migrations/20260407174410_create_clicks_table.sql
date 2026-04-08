-- +goose Up

CREATE TABLE clicks(
    id SERIAL PRIMARY KEY,
    url_id BIGINT NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    ip INET NOT NULL,
    user_agent TEXT,
    refer TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX clicks_url_id_created_at_idx ON clicks(url_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS clicks;