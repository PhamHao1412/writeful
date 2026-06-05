package service

import (
	"chat-service/internal/entity"
	"chat-service/internal/gateway/auth"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"chat-service/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type IConversationService interface {
	Create(ctx context.Context, userID string, req model.CreateConversationRequest) (*model.ConversationResponse, error, int)
	GetList(ctx context.Context, userID string, req model.GetConversationsRequest) ([]model.ConversationResponse, int64, error)
	GetByID(ctx context.Context, userID, conversationID string) (*model.ConversationResponse, error, int)
	AddParticipants(ctx context.Context, userID, conversationID string, req model.AddParticipantsRequest) error
	RemoveParticipant(ctx context.Context, userID, conversationID, participantID string) error
	MarkAsRead(ctx context.Context, userID string, req model.MarkAsReadRequest) error
	AutoRejoinParticipants(ctx context.Context, conversationID string, senderID string) error
}

type ConversationService struct {
	conversationRepo repository.IConversationRepository
	messageRepo      repository.IMessageRepository
	authClient       auth.Client
	hub              *Hub
}

func NewConversationService(
	conversationRepo repository.IConversationRepository,
	messageRepo repository.IMessageRepository,
	authClient auth.Client,
	hub *Hub,
) IConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		authClient:       authClient,
		hub:              hub,
	}
}

func (s *ConversationService) Create(ctx context.Context, userID string, req model.CreateConversationRequest) (*model.ConversationResponse, error, int) {
	// Validate request
	if req.Type == string(entity.ConversationTypeDirect) {
		if len(req.ParticipantIDs) != 1 {
			return nil, errors.New("direct conversation must have exactly 1 other participant"), http.StatusBadRequest
		}

		// Check if direct conversation already exists
		existing, err := s.conversationRepo.GetDirectConversation(ctx, userID, req.ParticipantIDs[0])
		if err == nil && existing != nil {
			// Kiểm tra user A đã left chưa
			hasLeft, err := s.conversationRepo.HasUserLeft(ctx, existing.ID, userID)
			if err == nil && hasLeft {
				// User A đã left → rejoin vào conversation cũ
				if err := s.conversationRepo.RejoinParticipant(ctx, existing.ID, userID); err != nil {
					return nil, err, http.StatusInternalServerError
				}
			}

			// Load lại conversation với participants mới
			conversation, err := s.conversationRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, err, http.StatusInternalServerError
			}

			return s.toConversationResponse(ctx, conversation, userID, nil), nil, http.StatusOK
		}
	}

	// Create conversation
	conversation := &entity.Conversation{
		Type:      entity.ConversationType(req.Type),
		Name:      req.Name,
		CreatedBy: userID,
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// Add participants (including creator)
	allParticipants := append([]string{userID}, req.ParticipantIDs...)
	if err := s.conversationRepo.AddParticipants(ctx, conversation.ID, allParticipants); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// Reload conversation with participants
	conversation, err := s.conversationRepo.GetByID(ctx, conversation.ID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return s.toConversationResponse(ctx, conversation, userID, nil), nil, http.StatusCreated
}

func (s *ConversationService) GetList(ctx context.Context, userID string, req model.GetConversationsRequest) ([]model.ConversationResponse, int64, error) {
	conversations, total, err := s.conversationRepo.GetByUserID(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	// Gather all unique user IDs across all conversations
	userIdsMap := make(map[string]bool)
	for _, conv := range conversations {
		for _, p := range conv.Participants {
			userIdsMap[p.UserID] = true
		}
	}

	uniqueUserIDs := make([]string, 0, len(userIdsMap))
	for id := range userIdsMap {
		uniqueUserIDs = append(uniqueUserIDs, id)
	}

	// Batch-fetch all user profiles from AuthService
	userMap := make(map[string]model.User)
	if len(uniqueUserIDs) > 0 {
		users, err, _ := s.authClient.GetListUser(model.GetUserRequest{
			IDs: uniqueUserIDs,
		})
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("[GetList] batch get list user failed err: %v", err))
		} else {
			for _, user := range users {
				userMap[user.ID] = user
			}
		}
	}

	responses := make([]model.ConversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		resp := s.toConversationResponse(ctx, &conv, userID, userMap)
		if resp != nil {
			responses = append(responses, *resp)
		}
	}

	return responses, total, nil
}

func (s *ConversationService) GetByID(ctx context.Context, userID, conversationID string) (*model.ConversationResponse, error, int) {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !isParticipant {
		return nil, errors.New("you are not a participant of this conversation"), http.StatusForbidden
	}

	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, err, http.StatusNotFound
	}

	return s.toConversationResponse(ctx, conversation, userID, nil), nil, http.StatusOK
}

func (s *ConversationService) AddParticipants(ctx context.Context, userID, conversationID string, req model.AddParticipantsRequest) error {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("you are not a participant of this conversation")
	}

	// Get conversation to check type
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return err
	}

	// Only group conversations can add participants
	if conversation.Type == entity.ConversationTypeDirect {
		return errors.New("cannot add participants to direct conversation")
	}

	return s.conversationRepo.AddParticipants(ctx, conversationID, req.UserIDs)
}

func (s *ConversationService) RemoveParticipant(ctx context.Context, userID, conversationID, participantID string) error {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("you are not a participant of this conversation")
	}

	// Users can only remove themselves or (if creator) remove others
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return err
	}

	if userID != participantID && conversation.CreatedBy != userID {
		return errors.New("you can only remove yourself or (if creator) remove others")
	}

	return s.conversationRepo.RemoveParticipant(ctx, conversationID, participantID)
}

func (s *ConversationService) MarkAsRead(ctx context.Context, userID string, req model.MarkAsReadRequest) error {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, req.ConversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("you are not a participant of this conversation")
	}

	err = s.conversationRepo.UpdateLastRead(ctx, req.ConversationID, userID)
	if err != nil {
		return err
	}

	// Fetch all participants to find other user IDs
	participants, err := s.conversationRepo.GetParticipants(ctx, req.ConversationID)
	if err == nil {
		var targetUserIDs []string
		for _, p := range participants {
			if p.UserID != userID && p.LeftAt == nil {
				targetUserIDs = append(targetUserIDs, p.UserID)
			}
		}

		// Broadcast "read" event via WebSocket to target users in background
		if len(targetUserIDs) > 0 && s.hub != nil {
			nowStr := time.Now().Format(time.RFC3339)
			s.hub.BroadcastToUsers(targetUserIDs, model.WSMessageTypeRead, map[string]interface{}{
				"conversation_id": req.ConversationID,
				"user_id":         userID,
				"read_at":         nowStr,
			})
		}
	}

	return nil
}

func (s *ConversationService) toConversationResponse(ctx context.Context, conv *entity.Conversation, userID string, userMap map[string]model.User) *model.ConversationResponse {
	participants := make([]model.ParticipantResponse, len(conv.Participants))

	if userMap == nil {
		userMap = make(map[string]model.User)
	}

	// Identify missing users from userMap
	var missingUserIDs []string
	for _, p := range conv.Participants {
		if _, exists := userMap[p.UserID]; !exists {
			missingUserIDs = append(missingUserIDs, p.UserID)
		}
	}

	// Fetch missing users from AuthService if any
	if len(missingUserIDs) > 0 {
		users, err, _ := s.authClient.GetListUser(model.GetUserRequest{
			IDs: missingUserIDs,
		})
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("[toConversationResponse] get list user failed err: %v", err))
		} else {
			for _, user := range users {
				userMap[user.ID] = user
			}
		}
	}

	for i, p := range conv.Participants {
		u, exists := userMap[p.UserID]
		if !exists {
			u = model.User{
				ID:          p.UserID,
				Username:    "unknown",
				DisplayName: "Unknown User",
			}
		}
		participants[i] = model.ParticipantResponse{
			ID:             p.ID,
			ConversationID: p.ConversationID,
			UserID:         p.UserID,
			JoinedAt:       p.JoinedAt,
			LastReadAt:     p.LastReadAt,
			User:           u,
		}
	}

	participant, err := s.conversationRepo.GetParticipant(ctx, conv.ID, userID)
	if err != nil {
		return nil
	}
	// Get last message
	var lastMessage *model.MessageResponse
	if lastMsg, err := s.messageRepo.GetLastMessage(ctx, conv.ID, participant.JoinedAt); err == nil && lastMsg != nil {
		lastMessage = &model.MessageResponse{
			ID:             lastMsg.ID,
			ConversationID: lastMsg.ConversationID,
			SenderID:       lastMsg.SenderID,
			Type:           string(lastMsg.Type),
			Content:        lastMsg.Content,
			MediaURL:       lastMsg.MediaURL,
			CreatedAt:      lastMsg.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      lastMsg.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Get unread count
	unreadCount, _ := s.conversationRepo.GetUnreadCount(ctx, conv.ID, userID)

	return &model.ConversationResponse{
		ID:            conv.ID,
		Type:          string(conv.Type),
		Name:          conv.Name,
		CreatedBy:     conv.CreatedBy,
		LastMessageAt: conv.LastMessageAt,
		Participants:  participants,
		LastMessage:   lastMessage,
		UnreadCount:   unreadCount,
		CreatedAt:     conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     conv.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *ConversationService) AutoRejoinParticipants(ctx context.Context, conversationID string, senderID string) error {
	participants, err := s.conversationRepo.GetAllParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	for _, p := range participants {
		if p.LeftAt != nil && p.UserID != senderID {
			if err := s.conversationRepo.RejoinParticipant(ctx, conversationID, p.UserID); err != nil {
				continue
			}
		}
	}

	return nil
}
