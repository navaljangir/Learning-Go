-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT todos_title_check CHECK (char_length(title) > 0)
);

-- Create indexes for better query performance
CREATE INDEX idx_todos_user_id ON todos(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_completed ON todos(completed) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_priority ON todos(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_due_date ON todos(due_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_todos_created_at ON todos(created_at DESC);
CREATE INDEX idx_todos_deleted_at ON todos(deleted_at);

-- Composite index for common queries
CREATE INDEX idx_todos_user_completed ON todos(user_id, completed) WHERE deleted_at IS NULL;

-- Create trigger for todos table
CREATE TRIGGER update_todos_updated_at BEFORE UPDATE ON todos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
