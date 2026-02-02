-- Drop trigger
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;

-- Drop indexes
DROP INDEX IF EXISTS idx_todos_user_completed;
DROP INDEX IF EXISTS idx_todos_deleted_at;
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_due_date;
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_completed;
DROP INDEX IF EXISTS idx_todos_user_id;

-- Drop table
DROP TABLE IF EXISTS todos;
