-- +migrate Up
ALTER TABLE users
  ADD COLUMN plan_expires_at TIMESTAMPTZ NULL;
