-- +goose Up

-- 1. Create shared musics library table
CREATE TABLE IF NOT EXISTS musics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    cover_url TEXT,
    genre VARCHAR(50) NOT NULL DEFAULT 'vpop',
    uploaded_by UUID NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_musics_genre ON musics(genre);

-- Insert initial premium curated seed tracks
INSERT INTO musics (title, artist, url, cover_url, genre) VALUES 
('Chúng Ta Của Tương Lai', 'Sơn Tùng M-TP', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3', 'https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=150', 'vpop'),
('Chúng Ta Của Hiện Tại', 'Sơn Tùng M-TP', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3', 'https://images.unsplash.com/photo-1498038432885-c6f3f1b912ee?w=150', 'vpop'),
('Muộn Rồi Mà Sao Còn', 'Sơn Tùng M-TP', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-5.mp3', 'https://images.unsplash.com/photo-1501386761578-eac5c94b800a?w=150', 'vpop'),
('Blinding Lights', 'The Weeknd', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3', 'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=150', 'pop'),
('Save Your Tears', 'The Weeknd', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-4.mp3', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=150', 'pop'),
('Starboy', 'The Weeknd ft. Daft Punk', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-6.mp3', 'https://images.unsplash.com/photo-1518609878373-06d740f60d8b?w=150', 'pop'),
('Sunset Lofi Study', 'Lofi Dreamer', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-8.mp3', 'https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?w=150', 'lofi'),
('Sunny Days Acoustic', 'Guitar Forest', 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-10.mp3', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=150', 'acoustic');

-- 2. Create stories table
CREATE TABLE IF NOT EXISTS stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'image',
    media_url TEXT NOT NULL,
    caption TEXT,
    audio_url TEXT,
    audio_title VARCHAR(255),
    audio_artist VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active'
);

CREATE INDEX IF NOT EXISTS idx_stories_user_id_expires ON stories(user_id, expires_at);

-- 3. Create story_views table for read status tracking
CREATE TABLE IF NOT EXISTS story_views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    viewer_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_story_viewer UNIQUE (story_id, viewer_id)
);

CREATE INDEX IF NOT EXISTS idx_story_views_story_id ON story_views(story_id);

-- +goose Down
DROP TABLE IF EXISTS story_views;
DROP TABLE IF EXISTS stories;
DROP TABLE IF EXISTS musics;
