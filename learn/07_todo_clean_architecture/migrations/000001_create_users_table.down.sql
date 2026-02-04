-- Drop indexes
DROP INDEX IF EXISTS idx_users_deleted_at ON users;
DROP INDEX IF EXISTS idx_users_email ON users;
DROP INDEX IF EXISTS idx_users_username ON users;

-- Drop table
DROP TABLE IF EXISTS users;
