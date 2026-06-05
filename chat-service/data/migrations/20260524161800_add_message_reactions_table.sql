-- +goose Up
-- Create message_reactions table
CREATE TABLE IF NOT EXISTS chat_service.message_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL,
    user_id UUID NOT NULL,
    emoji VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_message_reaction FOREIGN KEY (message_id) REFERENCES chat_service.messages(id) ON DELETE CASCADE,
    CONSTRAINT uq_message_user_reaction UNIQUE (message_id, user_id)
);

-- Optimize queries by message_id
CREATE INDEX IF NOT EXISTS idx_message_reactions_message_id ON chat_service.message_reactions(message_id);

-- +goose Down
DROP TABLE IF EXISTS chat_service.message_reactions;
