-- +migrate Up
CREATE TABLE apikeys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    key TEXT NOT NULL UNIQUE,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_apikeys_user_id ON apikeys(user_id);
