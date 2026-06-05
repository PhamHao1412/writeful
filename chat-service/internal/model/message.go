package model

type SendMessageRequest struct {
	ConversationID   string `json:"conversation_id" binding:"required"`
	Type             string `json:"type" binding:"required,oneof=text image file call"`
	Content          string `json:"content"`
	MediaURL         string `json:"media_url"`
	ReplyToMessageID string `json:"reply_to_message_id,omitempty"`
}

type GetMessagesRequest struct {
	ConversationID string `form:"conversation_id" binding:"required"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	PageSize       int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Before         string `form:"before"` // Message ID for pagination
}

type ToggleReactionRequest struct {
	Emoji string `json:"emoji" binding:"required"`
}

type ReactionResponse struct {
	ID        string `json:"id"`
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Emoji     string `json:"emoji"`
	CreatedAt string `json:"created_at"`
}

type MessageResponse struct {
	ID               string             `json:"id"`
	ConversationID   string             `json:"conversation_id"`
	SenderID         string             `json:"sender_id"`
	Type             string             `json:"type"`
	Content          string             `json:"content"`
	MediaURL         *string            `json:"media_url,omitempty"`
	ReplyToMessageID *string            `json:"reply_to_message_id,omitempty"`
	ReplyToMessage   *MessageResponse   `json:"reply_to_message,omitempty"`
	Reactions        []ReactionResponse `json:"reactions,omitempty"`
	CreatedAt        string             `json:"created_at"`
	UpdatedAt        string             `json:"updated_at"`
}

type MarkAsReadRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}

// WebSocket message types
type WSMessageType string

const (
	WSMessageTypeNewMessage   WSMessageType = "new_message"
	WSMessageTypeTyping       WSMessageType = "typing"
	WSMessageTypeRead         WSMessageType = "read"
	WSMessageTypeUserOnline   WSMessageType = "user_online"
	WSMessageTypeUserOffline  WSMessageType = "user_offline"
	WSMessageTypeCallInitiate WSMessageType = "call_initiate"
	WSMessageTypeCallReceive  WSMessageType = "call_receive"
	WSMessageTypeCallRinging  WSMessageType = "call_ringing"
	WSMessageTypeCallReject   WSMessageType = "call_reject"
	WSMessageTypeCallCancel   WSMessageType = "call_cancel"
	WSMessageTypeCallHangup   WSMessageType = "call_hangup"
	WSMessageTypeWebRTCOffer  WSMessageType = "webrtc_offer"
	WSMessageTypeWebRTCAnswer WSMessageType = "webrtc_answer"
	WSMessageTypeWebRTCICE    WSMessageType = "webrtc_ice_candidate"
	WSMessageTypeReaction     WSMessageType = "message_reaction"
)

type WSReactionPayload struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	Emoji          string `json:"emoji,omitempty"`
	Action         string `json:"action"` // "add", "update", "remove"
}

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type TypingPayload struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}
