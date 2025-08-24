-- +migrate Up
-- URL Shortener Database Schema
-- Complete schema with all features including password protection

-- Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL
);

-- API Keys table
CREATE TABLE apikeys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key TEXT NOT NULL UNIQUE,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Links table
CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    apikey_id BIGINT NULL REFERENCES apikeys(id) ON DELETE SET NULL,
    short_code TEXT NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    password TEXT NULL,
    click_count BIGINT NOT NULL DEFAULT 0,
    last_clicked_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_apikeys_key ON apikeys(key) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_links_short_code ON links(short_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_links_apikey_id ON links(apikey_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_apikeys_user_id ON apikeys(user_id) WHERE deleted_at IS NULL;

-- Add helpful comments
COMMENT ON TABLE users IS 'User accounts with optional password authentication';
COMMENT ON TABLE apikeys IS 'API keys for user authentication';
COMMENT ON TABLE links IS 'Shortened URLs with click tracking and optional password protection';
COMMENT ON COLUMN users.email IS 'Unique user email address';
COMMENT ON COLUMN users.password IS 'Optional password for user authentication';
COMMENT ON COLUMN apikeys.key IS '32-character hex API key';
COMMENT ON COLUMN links.apikey_id IS 'References the API key that owns this link (nullable)';
COMMENT ON COLUMN links.short_code IS 'Unique short code for the URL';
COMMENT ON COLUMN links.long_url IS 'Original long URL to redirect to';
COMMENT ON COLUMN links.password IS 'Optional password for link protection';
COMMENT ON COLUMN links.click_count IS 'Number of times this link was clicked';
COMMENT ON COLUMN links.last_clicked_at IS 'Timestamp of last click';
