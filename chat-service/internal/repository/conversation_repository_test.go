package repository

import (
	"chat-service/internal/entity"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	entity.SetSchemaName("")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(
		&entity.Conversation{},
		&entity.Participant{},
		&entity.Message{},
	)
	assert.NoError(t, err)

	return db
}

func TestConversationRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewConversationRepository(db)

	conversation := &entity.Conversation{
		Type:      entity.ConversationTypeDirect,
		CreatedBy: "user-1",
	}

	err := repo.Create(context.Background(), conversation)
	assert.NoError(t, err)
	assert.NotEmpty(t, conversation.ID)
}

func TestConversationRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewConversationRepository(db)

	// Create conversation
	conversation := &entity.Conversation{
		Type:      entity.ConversationTypeDirect,
		CreatedBy: "user-1",
	}
	err := repo.Create(context.Background(), conversation)
	assert.NoError(t, err)

	// Get conversation
	result, err := repo.GetByID(context.Background(), conversation.ID)
	assert.NoError(t, err)
	assert.Equal(t, conversation.ID, result.ID)
	assert.Equal(t, entity.ConversationTypeDirect, result.Type)
}

func TestConversationRepository_AddParticipants(t *testing.T) {
	db := setupTestDB(t)
	repo := NewConversationRepository(db)

	// Create conversation
	conversation := &entity.Conversation{
		Type:      entity.ConversationTypeGroup,
		Name:      "Test Group",
		CreatedBy: "user-1",
	}
	err := repo.Create(context.Background(), conversation)
	assert.NoError(t, err)

	// Add participants
	userIDs := []string{"user-1", "user-2", "user-3"}
	err = repo.AddParticipants(context.Background(), conversation.ID, userIDs)
	assert.NoError(t, err)

	// Verify participants
	participants, err := repo.GetParticipants(context.Background(), conversation.ID)
	assert.NoError(t, err)
	assert.Len(t, participants, 3)
}

func TestConversationRepository_IsParticipant(t *testing.T) {
	db := setupTestDB(t)
	repo := NewConversationRepository(db)

	// Create conversation and add participant
	conversation := &entity.Conversation{
		Type:      entity.ConversationTypeDirect,
		CreatedBy: "user-1",
	}
	err := repo.Create(context.Background(), conversation)
	assert.NoError(t, err)

	err = repo.AddParticipants(context.Background(), conversation.ID, []string{"user-1", "user-2"})
	assert.NoError(t, err)

	// Check if user is participant
	isParticipant, err := repo.IsParticipant(context.Background(), conversation.ID, "user-1")
	assert.NoError(t, err)
	assert.True(t, isParticipant)

	// Check non-participant
	isParticipant, err = repo.IsParticipant(context.Background(), conversation.ID, "user-3")
	assert.NoError(t, err)
	assert.False(t, isParticipant)
}
