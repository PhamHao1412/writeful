package model

type CreateConversationRequest struct {
	Type           string   `json:"type" binding:"required,oneof=direct group"`
	Name           string   `json:"name"`
	ParticipantIDs []string `json:"participant_ids" binding:"required,min=1"`
}

type GetConversationsRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type ConversationResponse struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Name          string                `json:"name,omitempty"`
	CreatedBy     string                `json:"created_by"`
	LastMessageAt *string               `json:"last_message_at,omitempty"`
	Participants  []ParticipantResponse `json:"participants,omitempty"`
	LastMessage   *MessageResponse      `json:"last_message,omitempty"`
	UnreadCount   int64                 `json:"unread_count"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

type ParticipantResponse struct {
	ID             string  `json:"id"`
	ConversationID string  `json:"conversation_id"`
	UserID         string  `json:"user_id"`
	JoinedAt       string  `json:"joined_at"`
	LastReadAt     *string `json:"last_read_at,omitempty"`
	User           User    `json:"user"`
}

type AddParticipantsRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
}
