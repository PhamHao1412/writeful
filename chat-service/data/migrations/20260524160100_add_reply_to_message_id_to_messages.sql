-- +goose Up
ALTER TABLE chat_service.messages
ADD COLUMN reply_to_message_id UUID REFERENCES chat_service.messages(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE chat_service.messages
DROP COLUMN reply_to_message_id;
