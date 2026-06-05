-- +goose Up
-- Modify messages type constraint to include 'call'
ALTER TABLE chat_service.messages DROP CONSTRAINT IF EXISTS messages_type_check;
ALTER TABLE chat_service.messages ADD CONSTRAINT messages_type_check CHECK (type IN ('text', 'image', 'file', 'call'));

-- +goose Down
-- Revert messages type constraint back to original
ALTER TABLE chat_service.messages DROP CONSTRAINT IF EXISTS messages_type_check;
ALTER TABLE chat_service.messages ADD CONSTRAINT messages_type_check CHECK (type IN ('text', 'image', 'file'));
