-- Таблица данных пользователя (PostgreSQL)
CREATE TABLE IF NOT EXISTS data (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    name TEXT NOT NULL,
    encrypted_data BYTEA NOT NULL,
    metadata TEXT,
    version BIGINT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_data_user_id ON data(user_id);
CREATE INDEX IF NOT EXISTS idx_data_type ON data(type);
CREATE INDEX IF NOT EXISTS idx_data_deleted_at ON data(deleted_at);
