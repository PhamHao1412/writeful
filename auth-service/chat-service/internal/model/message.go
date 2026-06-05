package model

type SendMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=text image file"`
	Content        string `json:"content"`
	MediaURL       string `json:"media_url"`
}

type GetMessagesRequest struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	PageSize       int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Before         string `form:"before"` // Message ID for pagination
}

type MessageResponse struct {
	ID             string  `json:"id"`
	ConversationID string  `json:"conversation_id"`
	SenderID       string  `json:"sender_id"`
	Type           string  `json:"type"`
	Content        string  `json:"content"`
	MediaURL       *string `json:"media_url,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type MarkAsReadRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}

// WebSocket message types
type WSMessageType string

const (
	WSMessageTypeNewMessage  WSMessageType = "new_message"
	WSMessageTypeTyping      WSMessageType = "typing"
	WSMessageTypeRead        WSMessageType = "read"
	WSMessageTypeUserOnline  WSMessageType = "user_online"
	WSMessageTypeUserOffline WSMessageType = "user_offline"
)

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type TypingPayload struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}
