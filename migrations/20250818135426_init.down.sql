-- +migrate Down
DROP INDEX IF EXISTS idx_user_longurl_unique;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;