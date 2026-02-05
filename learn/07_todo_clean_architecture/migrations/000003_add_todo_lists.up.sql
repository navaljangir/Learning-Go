-- Create todo_lists table
CREATE TABLE IF NOT EXISTS todo_lists (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_lists_user_id FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_user_lists (user_id, created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add list_id column to todos (nullable for global todos)
ALTER TABLE todos ADD COLUMN list_id CHAR(36) NULL AFTER user_id;

-- Add foreign key with CASCADE on delete (deletes todos with list)
ALTER TABLE todos ADD CONSTRAINT fk_todos_list_id
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE;

-- Add index for list queries
CREATE INDEX idx_todos_list_id ON todos(list_id);
