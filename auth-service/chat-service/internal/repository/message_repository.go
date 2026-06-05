package repository

import (
	"chat-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type IMessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	GetByID(ctx context.Context, id string) (*entity.Message, error)
	GetByConversationID(ctx context.Context, conversationID string, page, pageSize int, beforeID string) ([]entity.Message, int64, error)
	Update(ctx context.Context, message *entity.Message) error
	Delete(ctx context.Context, id string) error
	GetLastMessage(ctx context.Context, conversationID string) (*entity.Message, error)
}

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) IMessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *entity.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepository) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	var message entity.Message
	err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) GetByConversationID(ctx context.Context, conversationID string, page, pageSize int, beforeID string) ([]entity.Message, int64, error) {
	var messages []entity.Message
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("conversation_id = ?", conversationID)

	// Cursor-based pagination using beforeID
	if beforeID != "" {
		var beforeMessage entity.Message
		if err := r.db.WithContext(ctx).First(&beforeMessage, "id = ?", beforeID).Error; err == nil {
			query = query.Where("created_at < ?", beforeMessage.CreatedAt)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error

	return messages, total, err
}

func (r *MessageRepository) Update(ctx context.Context, message *entity.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *MessageRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Message{}, "id = ?", id).Error
}

func (r *MessageRepository) GetLastMessage(ctx context.Context, conversationID string) (*entity.Message, error) {
	var message entity.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}
