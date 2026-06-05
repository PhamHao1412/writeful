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

type IMessageService interface {
	Send(ctx context.Context, userID string, req model.SendMessageRequest) (*model.MessageResponse, error, int)
	GetMessages(ctx context.Context, userID string, req model.GetMessagesRequest) ([]model.MessageResponse, int64, error, int)
	Delete(ctx context.Context, userID, messageID string) error
}

type MessageService struct {
	messageRepo      repository.IMessageRepository
	conversationRepo repository.IConversationRepository
}

func NewMessageService(
	messageRepo repository.IMessageRepository,
	conversationRepo repository.IConversationRepository,
) IMessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
	}
}

func (s *MessageService) Send(ctx context.Context, userID string, req model.SendMessageRequest) (*model.MessageResponse, error, int) {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, req.ConversationID, userID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !isParticipant {
		return nil, errors.New("you are not a participant of this conversation"), http.StatusForbidden
	}

	// Validate message content
	if req.Type == string(entity.MessageTypeText) && req.Content == "" {
		return nil, errors.New("text message must have content"), http.StatusBadRequest
	}
	if (req.Type == string(entity.MessageTypeImage) || req.Type == string(entity.MessageTypeFile)) && req.MediaURL == "" {
		return nil, errors.New("media message must have media URL"), http.StatusBadRequest
	}

	// Create message
	var mediaURL *string
	if req.MediaURL != "" {
		mediaURL = &req.MediaURL
	}

	message := &entity.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Type:           entity.MessageType(req.Type),
		Content:        req.Content,
		MediaURL:       mediaURL,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// Update conversation's last_message_at
	conversation, err := s.conversationRepo.GetByID(ctx, req.ConversationID)
	if err == nil {
		now := time.Now().Format(time.RFC3339)
		conversation.LastMessageAt = &now
		s.conversationRepo.Update(ctx, conversation)
	}

	return &model.MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Type:           string(message.Type),
		Content:        message.Content,
		MediaURL:       message.MediaURL,
		CreatedAt:      message.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      message.UpdatedAt.Format(time.RFC3339),
	}, nil, http.StatusCreated
}

func (s *MessageService) GetMessages(ctx context.Context, userID string, req model.GetMessagesRequest) ([]model.MessageResponse, int64, error, int) {
	// Check if user is participant
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, req.ConversationID, userID)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}
	if !isParticipant {
		return nil, 0, errors.New("you are not a participant of this conversation"), http.StatusForbidden
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}

	messages, total, err := s.messageRepo.GetByConversationID(ctx, req.ConversationID, req.Page, req.PageSize, req.Before)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}

	responses := make([]model.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = model.MessageResponse{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			SenderID:       msg.SenderID,
			Type:           string(msg.Type),
			Content:        msg.Content,
			MediaURL:       msg.MediaURL,
			CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      msg.UpdatedAt.Format(time.RFC3339),
		}
	}

	return responses, total, nil, http.StatusOK
}

func (s *MessageService) Delete(ctx context.Context, userID, messageID string) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Only sender can delete message
	if message.SenderID != userID {
		return errors.New("you can only delete your own messages")
	}

	return s.messageRepo.Delete(ctx, messageID)
}
