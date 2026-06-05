-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE media_service.images
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    url TEXT NOT NULL,
    format TEXT NOT NULL,
    type VARCHAR(20) NOT NULL,

    provider VARCHAR(20) DEFAULT 'cloudinary',
    public_id TEXT,

    mime_type TEXT,
    file_size BIGINT,
    width INT,
    height INT,
    uploaded_by TEXT

);
CREATE TABLE videos
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type       text,
    url        text,
    format     text,
    duration   INTEGER,
    size       BIGINT,
    thumbnail_url text,

    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- +goose Down
DROP TABLE images;
DROP TABLE videos;
