-- +goose Up
-- Create conversations table
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS chat_service.conversations (
                                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                    type VARCHAR(20) NOT NULL CHECK (type IN ('direct', 'group')),
                                                    name VARCHAR(255),
                                                    created_by UUID NOT NULL,
                                                    last_message_at TIMESTAMP,
                                                    is_started boolean default false,
                                                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                    deleted_at TIMESTAMP
);

-- Create participants table
CREATE TABLE IF NOT EXISTS chat_service.participants (
                                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                   conversation_id UUID NOT NULL,
                                                   user_id UUID NOT NULL,
                                                   joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                   last_read_at TIMESTAMP,
                                                   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                   deleted_at TIMESTAMP,
                                                   left_at TIMESTAMP,
                                                   joined_at TIMESTAMP,
                                                   CONSTRAINT fk_conversation FOREIGN KEY (conversation_id) REFERENCES chat_service.conversations(id) ON DELETE CASCADE
);

-- Create messages table
CREATE TABLE IF NOT EXISTS chat_service.messages (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               conversation_id UUID NOT NULL,
                                               sender_id UUID NOT NULL,
                                               type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'image', 'file')),
                                               content TEXT,
                                               media_url TEXT,
                                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                               deleted_at TIMESTAMP,
                                               CONSTRAINT fk_conversation_msg FOREIGN KEY (conversation_id) REFERENCES chat_service.conversations(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS chat_service.messages;
DROP TABLE IF EXISTS chat_service.participants;
DROP TABLE IF EXISTS chat_service.conversations;