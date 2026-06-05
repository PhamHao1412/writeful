package entity

import "fmt"

type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct"
	ConversationTypeGroup  ConversationType = "group"
)

type Conversation struct {
	BaseEntity
	Type          ConversationType `gorm:"column:type;type:varchar(20);not null" json:"type"`
	Name          string           `gorm:"column:name;type:varchar(255)" json:"name,omitempty"`
	CreatedBy     string           `gorm:"column:created_by;not null" json:"created_by"`
	LastMessageAt *string          `gorm:"column:last_message_at" json:"last_message_at,omitempty"`
	IsStared      bool             `gorm:"column:is_started" json:"is_started"`
	// Relations
	Participants []Participant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	Messages     []Message     `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return fmt.Sprintf("%vconversations", SchemaName())
}
func (Conversation) TableNameAlias(alias string) string {
	return SchemaName() + "conversations " + alias
}
