package service

import (
	"chat-service/internal/entity"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) GetUserProfile(req model.GetUserRequest) (*model.User, error, int) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1), args.Int(2)
	}
	return args.Get(0).(*model.User), args.Error(1), args.Int(2)
}

func (m *MockAuthClient) GetListUser(req model.GetUserRequest) ([]model.User, error, int) {
	args := m.Called(req)
	return args.Get(0).([]model.User), args.Error(1), args.Int(2)
}

func TestConversationService_GetList_PanicAndN1Fix(t *testing.T) {
	entity.SetSchemaName("")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&entity.Conversation{},
		&entity.Participant{},
		&entity.Message{},
	)
	assert.NoError(t, err)

	convRepo := repository.NewConversationRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	mockAuth := new(MockAuthClient)

	svc := NewConversationService(convRepo, msgRepo, mockAuth)

	// Create test conversation
	conv1 := &entity.Conversation{
		Type:      entity.ConversationTypeDirect,
		CreatedBy: "user-1",
		IsStarted: true,
	}
	err = convRepo.Create(context.Background(), conv1)
	assert.NoError(t, err)

	err = convRepo.AddParticipants(context.Background(), conv1.ID, []string{"user-1", "user-2"})
	assert.NoError(t, err)

	// Case 1: AuthService returns users successfully. GetList should batch call GetListUser.
	mockAuth.On("GetListUser", model.GetUserRequest{
		IDs: []string{"user-1", "user-2"},
	}).Return([]model.User{
		{ID: "user-1", Username: "user1", DisplayName: "User One"},
		{ID: "user-2", Username: "user2", DisplayName: "User Two"},
	}, nil, 200).Once()

	res, total, err := svc.GetList(context.Background(), "user-1", model.GetConversationsRequest{Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, res, 1)
	assert.Len(t, res[0].Participants, 2)
	// Sort order is not guaranteed by uniqueIDs map key iteration, let's identify users by ID
	user1Found := false
	user2Found := false
	for _, p := range res[0].Participants {
		if p.UserID == "user-1" {
			assert.Equal(t, "User One", p.User.DisplayName)
			user1Found = true
		} else if p.UserID == "user-2" {
			assert.Equal(t, "User Two", p.User.DisplayName)
			user2Found = true
		}
	}
	assert.True(t, user1Found)
	assert.True(t, user2Found)

	// Case 2: AuthService returns empty/incomplete list. It should gracefully fallback to placeholders instead of panicking.
	mockAuth.On("GetListUser", model.GetUserRequest{
		IDs: []string{"user-1", "user-2"},
	}).Return([]model.User{}, nil, 200).Once() // empty user response

	res, total, err = svc.GetList(context.Background(), "user-1", model.GetConversationsRequest{Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, res, 1)
	assert.Len(t, res[0].Participants, 2)
	assert.Equal(t, "Unknown User", res[0].Participants[0].User.DisplayName)
	assert.Equal(t, "unknown", res[0].Participants[0].User.Username)

	mockAuth.AssertExpectations(t)
}
