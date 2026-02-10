-- Таблица пользователей (SQLite)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
