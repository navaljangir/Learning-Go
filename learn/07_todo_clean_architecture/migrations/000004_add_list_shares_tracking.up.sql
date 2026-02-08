-- Create list_shares table to track sharing history
-- This is optional but useful for auditing and analytics
CREATE TABLE IF NOT EXISTS list_shares (
    id CHAR(36) PRIMARY KEY,
    source_list_id CHAR(36) NOT NULL,
    source_user_id CHAR(36) NOT NULL,
    target_list_id CHAR(36) NOT NULL,
    target_user_id CHAR(36) NOT NULL,
    shared_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    CONSTRAINT fk_shares_source_user FOREIGN KEY (source_user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_shares_target_user FOREIGN KEY (target_user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_shares_target_list FOREIGN KEY (target_list_id)
        REFERENCES todo_lists(id) ON DELETE CASCADE,

    -- Indexes for efficient queries
    INDEX idx_source_user_shares (source_user_id, shared_at),
    INDEX idx_target_user_shares (target_user_id, shared_at),
    INDEX idx_source_list_shares (source_list_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
