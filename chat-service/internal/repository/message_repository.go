package repository

import (
	"chat-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type IMessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	GetByID(ctx context.Context, id string) (*entity.Message, error)
	GetByConversationID(ctx context.Context, conversationID string, page, pageSize int, beforeID, afterTime string) ([]entity.Message, int64, error)
	Update(ctx context.Context, message *entity.Message) error
	Delete(ctx context.Context, id string) error
	GetLastMessage(ctx context.Context, conversationID, afterTime string) (*entity.Message, error)
	Count(ctx context.Context, conversationID string) (int64, error)
	ToggleReaction(ctx context.Context, reaction *entity.MessageReaction) (string, error)
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
	err := r.db.WithContext(ctx).
		Preload("ReplyToMessage").
		Preload("Reactions").
		First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) GetByConversationID(ctx context.Context, conversationID string, page, pageSize int, beforeID string, afterTime string) ([]entity.Message, int64, error) {
	var messages []entity.Message
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("conversation_id = ? AND created_at >= ?", conversationID, afterTime)

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
		Preload("ReplyToMessage").
		Preload("Reactions").
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

func (r *MessageRepository) GetLastMessage(ctx context.Context, conversationID, afterTime string) (*entity.Message, error) {
	var message entity.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND created_at >= ? ", conversationID, afterTime).
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

func (r *MessageRepository) Count(ctx context.Context, conversationID string) (int64, error) {
	db := r.db.WithContext(ctx)
	var count int64
	err := db.Model(&entity.Message{}).Where("conversation_id = ?", conversationID).Count(&count).Error
	return count, err
}

func (r *MessageRepository) ToggleReaction(ctx context.Context, reaction *entity.MessageReaction) (string, error) {
	var existing entity.MessageReaction
	err := r.db.WithContext(ctx).
		Where("message_id = ? AND user_id = ?", reaction.MessageID, reaction.UserID).
		First(&existing).Error

	if err == nil {
		// Reaction already exists
		if existing.Emoji == reaction.Emoji {
			// Emojis match, delete (remove) the reaction
			err = r.db.WithContext(ctx).Delete(&existing).Error
			return "remove", err
		} else {
			// Emojis differ, update the reaction to the new emoji
			existing.Emoji = reaction.Emoji
			err = r.db.WithContext(ctx).Save(&existing).Error
			return "update", err
		}
	} else if err == gorm.ErrRecordNotFound {
		// Reaction does not exist, insert a new reaction
		err = r.db.WithContext(ctx).Create(reaction).Error
		return "add", err
	}
	return "", err
}
