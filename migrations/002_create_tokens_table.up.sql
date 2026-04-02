-- +migrate Up
-- Create tokens table for storing refresh tokens
CREATE TABLE IF NOT EXISTS tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    token_type TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);

-- Create index on token_type for filtering
CREATE INDEX IF NOT EXISTS idx_tokens_token_type ON tokens(token_type);

-- Create token_blacklist table for revoked tokens
CREATE TABLE IF NOT EXISTS token_blacklist (
    id TEXT PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on expires_at for cleanup
CREATE INDEX IF NOT EXISTS idx_blacklist_expires_at ON token_blacklist(expires_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_blacklist_expires_at;
DROP TABLE IF EXISTS token_blacklist;
DROP INDEX IF EXISTS idx_tokens_token_type;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP TABLE IF EXISTS tokens;
