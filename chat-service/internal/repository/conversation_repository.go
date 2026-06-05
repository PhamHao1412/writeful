package repository

import (
	"chat-service/internal/entity"
	"context"
	"time"

	"gorm.io/gorm"
)

type IConversationRepository interface {
	Create(ctx context.Context, conversation *entity.Conversation) error
	GetByID(ctx context.Context, id string) (*entity.Conversation, error)
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]entity.Conversation, int64, error)
	GetDirectConversation(ctx context.Context, user1ID, user2ID string) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation) error
	Delete(ctx context.Context, id string) error
	StartConversation(ctx context.Context, id string) error

	// Participant operations
	AddParticipants(ctx context.Context, conversationID string, userIDs []string) error
	RemoveParticipant(ctx context.Context, conversationID, userID string) error
	GetParticipants(ctx context.Context, conversationID string) ([]entity.Participant, error)
	IsParticipant(ctx context.Context, conversationID, userID string) (bool, error)
	GetParticipant(ctx context.Context, conversationID, userID string) (*entity.Participant, error)
	UpdateLastRead(ctx context.Context, conversationID, userID string) error
	GetUnreadCount(ctx context.Context, conversationID, userID string) (int64, error)
	HasUserLeft(ctx context.Context, conversationID, userID string) (bool, error)
	RejoinParticipant(ctx context.Context, conversationID, userID string) error
	GetAllParticipants(ctx context.Context, conversationID string) ([]entity.Participant, error)
}

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) IConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, conversation *entity.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*entity.Conversation, error) {
	var conversation entity.Conversation
	err := r.db.WithContext(ctx).
		Preload("Participants").
		First(&conversation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]entity.Conversation, int64, error) {
	var conversations []entity.Conversation
	var total int64

	offset := (page - 1) * pageSize

	// Get conversation IDs where user is a participant
	subQuery := r.db.Model(&entity.Participant{}).
		Select("conversation_id").
		Where("user_id = ? AND left_at IS NULL", userID)

	query := r.db.WithContext(ctx).
		Model(&entity.Conversation{}).
		Where("id IN (?)", subQuery).Where("created_by = ? OR is_started = ?", userID, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Participants").
		Order("last_message_at DESC NULLS LAST, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&conversations).Error

	return conversations, total, err
}

func (r *ConversationRepository) GetDirectConversation(ctx context.Context, user1ID, user2ID string) (*entity.Conversation, error) {
	var conversation entity.Conversation

	// Find direct conversation between two users
	err := r.db.WithContext(ctx).
		Joins("JOIN "+entity.SchemaName()+"participants p1 ON p1.conversation_id = "+entity.SchemaName()+"conversations.id").
		Joins("JOIN "+entity.SchemaName()+"participants p2 ON p2.conversation_id = "+entity.SchemaName()+"conversations.id").
		Where(entity.SchemaName()+"conversations.type = ?", entity.ConversationTypeDirect).
		Where("p1.user_id = ? AND p2.user_id = ?", user1ID, user2ID).
		Preload("Participants").
		First(&conversation).Error

	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) Update(ctx context.Context, conversation *entity.Conversation) error {
	return r.db.WithContext(ctx).Save(conversation).Error
}

func (r *ConversationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Conversation{}, "id = ?", id).Error
}

func (r *ConversationRepository) AddParticipants(ctx context.Context, conversationID string, userIDs []string) error {
	now := time.Now().Format(time.RFC3339)
	participants := make([]entity.Participant, len(userIDs))

	for i, userID := range userIDs {
		participants[i] = entity.Participant{
			ConversationID: conversationID,
			UserID:         userID,
			JoinedAt:       now,
		}
	}

	return r.db.WithContext(ctx).Create(&participants).Error
}

func (r *ConversationRepository) RemoveParticipant(ctx context.Context, conversationID, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Participant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		Update("left_at", now).Error
}

func (r *ConversationRepository) GetParticipants(ctx context.Context, conversationID string) ([]entity.Participant, error) {
	var participants []entity.Participant
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Find(&participants).Error
	return participants, err
}

func (r *ConversationRepository) IsParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Participant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *ConversationRepository) GetParticipant(ctx context.Context, conversationID, userID string) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		First(&participant).Error
	return &participant, err
}

func (r *ConversationRepository) UpdateLastRead(ctx context.Context, conversationID, userID string) error {
	now := time.Now().Format(time.RFC3339)
	return r.db.WithContext(ctx).
		Model(&entity.Participant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", now).Error
}

func (r *ConversationRepository) StartConversation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entity.Conversation{}).Where("id = ?", id).Update("is_started", true).Error
}

func (r *ConversationRepository) GetUnreadCount(ctx context.Context, conversationID, userID string) (int64, error) {
	var participant entity.Participant
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		First(&participant).Error
	if err != nil {
		return 0, err
	}

	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("conversation_id = ?", conversationID).
		Where("sender_id != ?", userID) // Exclude messages sent by current user

	if participant.LastReadAt != nil {
		query = query.Where("created_at > ?", *participant.LastReadAt)
	} else {
		// If never read, count messages after joined_at
		if participant.JoinedAt != "" {
			query = query.Where("created_at > ?", participant.JoinedAt)
		}
	}

	err = query.Count(&count).Error
	return count, err
}

func (r *ConversationRepository) HasUserLeft(ctx context.Context, conversationID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Participant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NOT NULL", conversationID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *ConversationRepository) RejoinParticipant(ctx context.Context, conversationID, userID string) error {
	now := time.Now().Format(time.RFC3339)
	return r.db.WithContext(ctx).
		Model(&entity.Participant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NOT NULL", conversationID, userID).
		Updates(map[string]interface{}{
			"left_at":      nil, // Clear left_at
			"joined_at":    now, // Update joined_at mới
			"last_read_at": nil, // Reset last_read_at
		}).Error
}

func (r *ConversationRepository) GetAllParticipants(ctx context.Context, conversationID string) ([]entity.Participant, error) {
	var participants []entity.Participant
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Find(&participants).Error
	return participants, err
}
