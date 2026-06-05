-- +goose Up
ALTER TABLE media_service.images ADD COLUMN IF NOT EXISTS file_hash VARCHAR(64);
CREATE INDEX IF NOT EXISTS idx_images_file_hash ON media_service.images(file_hash);

-- +goose Down
DROP INDEX IF EXISTS media_service.idx_images_file_hash;
ALTER TABLE media_service.images DROP COLUMN IF EXISTS file_hash;
