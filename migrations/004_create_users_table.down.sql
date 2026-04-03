-- +migrate Down
-- Drop users table and related constraints/indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_status;

-- Drop the table
DROP TABLE IF EXISTS users;
