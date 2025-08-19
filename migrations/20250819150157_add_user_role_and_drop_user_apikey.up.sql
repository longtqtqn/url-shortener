-- +migrate Up
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

-- Drop legacy column if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='users' AND column_name='apikey'
    ) THEN
        ALTER TABLE users DROP COLUMN apikey;
    END IF;
END$$;
