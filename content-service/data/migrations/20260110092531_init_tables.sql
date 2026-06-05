-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- POSTS
-- =========================
CREATE TABLE posts
(
    id             UUID PRIMARY KEY     DEFAULT gen_random_uuid(),

    user_id        UUID        NOT NULL,
    title          TEXT        NOT NULL,
    slug           TEXT        NOT NULL UNIQUE,
    visibility     VARCHAR(20) NOT NULL DEFAULT 'public',
    excerpt        TEXT,

    status         VARCHAR(20) NOT NULL DEFAULT 'draft',
    published_at   TIMESTAMP,

    cover_image_url TEXT,

    created_at     TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_status ON posts (status);
CREATE INDEX idx_posts_published_at ON posts (published_at);

-- =========================
-- POST VERSIONS (HISTORY)
-- =========================
CREATE TABLE post_versions
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id    UUID NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    version    INT  NOT NULL,

    content    TEXT NOT NULL,
    created_by UUID NOT NULL,

    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    UNIQUE (post_id, version)
);

CREATE INDEX idx_post_versions_post_id ON post_versions (post_id);
CREATE TABLE media
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP,

    url         TEXT        NOT NULL,
    format      TEXT        NOT NULL,
    type        VARCHAR(20) NOT NULL,

    provider    VARCHAR(20)      DEFAULT 'cloudinary',
    public_id   TEXT,

    mime_type   TEXT,
    file_size   BIGINT,
    width       INT,
    height      INT,

    uploaded_by UUID NOT NULL

);

CREATE INDEX idx_media_uploaded_by ON media (uploaded_by);
CREATE INDEX idx_media_type ON media (type);
CREATE INDEX idx_media_deleted_at ON media (deleted_at);
CREATE INDEX idx_media_public_id ON media (public_id);


CREATE TABLE post_media
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id       UUID NOT NULL REFERENCES posts (id),
    media_id      UUID NOT NULL REFERENCES media (id),
    display_order INT              DEFAULT 0,


    created_at    TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP
);

CREATE INDEX idx_post_media_post_id ON post_media (post_id);
CREATE INDEX idx_post_media_media_id ON post_media (media_id);
CREATE INDEX idx_post_media_deleted_at ON post_media (deleted_at);
-- =========================
-- TAGS
-- =========================
CREATE TABLE tags
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);

-- =========================
-- POST TAGS
-- =========================
CREATE TABLE post_tags
(
    post_id UUID NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- =========================
-- POST STATS (OPTIONAL BUT USEFUL)
-- =========================
CREATE TABLE post_stats
(
    post_id       UUID PRIMARY KEY REFERENCES posts (id) ON DELETE CASCADE,
    view_count    BIGINT DEFAULT 0,
    like_count    BIGINT DEFAULT 0,
    comment_count BIGINT DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS post_stats;
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS post_versions;
DROP TABLE IF EXISTS posts;
