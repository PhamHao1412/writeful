package entity

import "fmt"

type Participant struct {
	BaseEntity
	ConversationID string  `gorm:"column:conversation_id;not null;index:idx_conversation_user" json:"conversation_id"`
	UserID         string  `gorm:"column:user_id;not null;index:idx_conversation_user" json:"user_id"`
	JoinedAt       string  `gorm:"column:joined_at" json:"joined_at"`
	LastReadAt     *string `gorm:"column:last_read_at" json:"last_read_at,omitempty"`

	// Relations
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

func (Participant) TableName() string {
	return fmt.Sprintf("%vparticipants", SchemaName())
}
