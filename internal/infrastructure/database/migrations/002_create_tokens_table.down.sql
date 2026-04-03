-- +migrate Down
DROP INDEX IF EXISTS idx_blacklist_expires_at;
DROP TABLE IF EXISTS token_blacklist;
DROP INDEX IF EXISTS idx_tokens_token_type;
DROP INDEX IF EXISTS idx_tokens_user_id;
DROP TABLE IF EXISTS tokens;
