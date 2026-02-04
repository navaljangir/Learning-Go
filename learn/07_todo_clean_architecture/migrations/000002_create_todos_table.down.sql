-- Drop indexes
DROP INDEX IF EXISTS idx_todos_user_completed ON todos;
DROP INDEX IF EXISTS idx_todos_deleted_at ON todos;
DROP INDEX IF EXISTS idx_todos_created_at ON todos;
DROP INDEX IF EXISTS idx_todos_due_date ON todos;
DROP INDEX IF EXISTS idx_todos_priority ON todos;
DROP INDEX IF EXISTS idx_todos_completed ON todos;
DROP INDEX IF EXISTS idx_todos_user_id ON todos;

-- Drop table
DROP TABLE IF EXISTS todos;
