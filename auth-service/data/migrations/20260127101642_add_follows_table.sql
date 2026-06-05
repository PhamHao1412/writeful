-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
SELECT 'up SQL query';
CREATE TABLE auth_service.follows (
                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                follower_id TEXT NOT NULL,
                                following_id TEXT NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                deleted_at TIMESTAMP
);
CREATE INDEX idx_user_follower_id ON auth_service.follows (follower_id);
CREATE INDEX idx_user_following_id ON auth_service.follows (following_id);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS refresh_tokens CASCADE;
