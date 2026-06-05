package entity

import "fmt"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	BaseEntity
	ConversationID string      `gorm:"column:conversation_id;not null;index" json:"conversation_id"`
	SenderID       string      `gorm:"column:sender_id;not null" json:"sender_id"`
	Type           MessageType `gorm:"column:type;type:varchar(20);not null" json:"type"`
	Content        string      `gorm:"column:content;type:text" json:"content"`
	MediaURL       *string     `gorm:"column:media_url" json:"media_url,omitempty"`

	// Relations
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

func (Message) TableName() string {
	return fmt.Sprintf("%vmessages", SchemaName())
}
