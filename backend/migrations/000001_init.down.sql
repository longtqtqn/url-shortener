-- +migrate Down
-- Drop all tables in reverse order (due to foreign key constraints)

-- Drop links table first (references apikeys)
DROP TABLE IF EXISTS links;

-- Drop apikeys table (references users)
DROP TABLE IF EXISTS apikeys;

-- Drop users table last
DROP TABLE IF EXISTS users;
