-- Таблица данных пользователя (SQLite)
CREATE TABLE IF NOT EXISTS data (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    encrypted_data BLOB NOT NULL,
    metadata TEXT,
    version INTEGER DEFAULT 1,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_data_user_id ON data(user_id);
CREATE INDEX IF NOT EXISTS idx_data_type ON data(type);
CREATE INDEX IF NOT EXISTS idx_data_deleted_at ON data(deleted_at);
