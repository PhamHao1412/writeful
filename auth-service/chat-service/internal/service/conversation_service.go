package service

import (
	"chat-service/internal/entity"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"context"
	"errors"
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
}

type ConversationService struct {
	conversationRepo repository.IConversationRepository
	messageRepo      repository.IMessageRepository
}

func NewConversationService(
	conversationRepo repository.IConversationRepository,
	messageRepo repository.IMessageRepository,
) IConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
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
			return s.toConversationResponse(ctx, existing, userID), nil, http.StatusOK
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

	return s.toConversationResponse(ctx, conversation, userID), nil, http.StatusCreated
}

func (s *ConversationService) GetList(ctx context.Context, userID string, req model.GetConversationsRequest) ([]model.ConversationResponse, int64, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	conversations, total, err := s.conversationRepo.GetByUserID(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		responses[i] = *s.toConversationResponse(ctx, &conv, userID)
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

	return s.toConversationResponse(ctx, conversation, userID), nil, http.StatusOK
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

	return s.conversationRepo.UpdateLastRead(ctx, req.ConversationID, userID)
}

func (s *ConversationService) toConversationResponse(ctx context.Context, conv *entity.Conversation, userID string) *model.ConversationResponse {
	participants := make([]model.ParticipantResponse, len(conv.Participants))
	for i, p := range conv.Participants {
		participants[i] = model.ParticipantResponse{
			ID:             p.ID,
			ConversationID: p.ConversationID,
			UserID:         p.UserID,
			JoinedAt:       p.JoinedAt,
			LastReadAt:     p.LastReadAt,
		}
	}

	// Get last message
	var lastMessage *model.MessageResponse
	if lastMsg, err := s.messageRepo.GetLastMessage(ctx, conv.ID); err == nil && lastMsg != nil {
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
