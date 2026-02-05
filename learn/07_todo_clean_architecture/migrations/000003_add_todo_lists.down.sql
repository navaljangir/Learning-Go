-- Remove foreign key and index
ALTER TABLE todos DROP FOREIGN KEY fk_todos_list_id;
DROP INDEX idx_todos_list_id ON todos;

-- Remove list_id column
ALTER TABLE todos DROP COLUMN list_id;

-- Drop todo_lists table
DROP TABLE IF EXISTS todo_lists;
