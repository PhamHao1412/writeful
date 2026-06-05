package entity

import "fmt"

type MessageReaction struct {
	BaseEntity
	MessageID string `gorm:"column:message_id;not null;index" json:"message_id"`
	UserID    string `gorm:"column:user_id;not null" json:"user_id"`
	Emoji     string `gorm:"column:emoji;type:varchar(50);not null" json:"emoji"`

	// Relations
	Message *Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}

func (MessageReaction) TableName() string {
	return fmt.Sprintf("%vmessage_reactions", SchemaName())
}
