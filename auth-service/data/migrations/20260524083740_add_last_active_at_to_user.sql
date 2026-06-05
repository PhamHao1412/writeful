-- +goose Up
ALTER TABLE users ADD COLUMN last_active_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE users DROP COLUMN last_active_at;
