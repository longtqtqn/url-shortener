-- +migrate Down
DROP INDEX IF EXISTS idx_apikeys_user_id;
DROP TABLE IF EXISTS apikeys;
