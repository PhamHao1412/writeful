package service

import (
	"chat-service/internal/entity"
	"chat-service/internal/model"
	"chat-service/internal/repository"
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type IMessageService interface {
	Send(ctx context.Context, userID string, req model.SendMessageRequest) (*model.MessageResponse, error, int)
	GetMessages(ctx context.Context, userID string, req model.GetMessagesRequest) ([]model.MessageResponse, int64, error, int)
	Delete(ctx context.Context, userID, messageID string) error
	ToggleReaction(ctx context.Context, userID, messageID string, req model.ToggleReactionRequest) (*model.ReactionResponse, string, error, int)
	GetByID(ctx context.Context, id string) (*model.MessageResponse, error)
}

type MessageService struct {
	messageRepo      repository.IMessageRepository
	conversationRepo repository.IConversationRepository
	conversationSvc  IConversationService
}

func NewMessageService(
	messageRepo repository.IMessageRepository,
	conversationRepo repository.IConversationRepository,
	conversationSvc IConversationService,
) IMessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		conversationSvc:  conversationSvc,
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

	if err := s.conversationSvc.AutoRejoinParticipants(ctx, req.ConversationID, userID); err != nil {
		log.Printf("Warning: Failed to auto rejoin participants: %v", err)
	}

	// start the conversation if it's a new conversation
	count, err := s.messageRepo.Count(ctx, req.ConversationID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if count == 0 {
		err = s.conversationRepo.StartConversation(ctx, req.ConversationID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
	}

	// Create message
	var mediaURL *string
	if req.MediaURL != "" {
		mediaURL = &req.MediaURL
	}

	var replyToMessageID *string
	if req.ReplyToMessageID != "" {
		replyToMessageID = &req.ReplyToMessageID
	}

	message := &entity.Message{
		ConversationID:   req.ConversationID,
		SenderID:         userID,
		Type:             entity.MessageType(req.Type),
		Content:          req.Content,
		MediaURL:         mediaURL,
		ReplyToMessageID: replyToMessageID,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// Load preloaded message to get full parent info if it was a reply
	preloaded, err := s.messageRepo.GetByID(ctx, message.ID)
	if err == nil && preloaded != nil {
		message = preloaded
	}

	// Update conversation's last_message_at
	conversation, err := s.conversationRepo.GetByID(ctx, req.ConversationID)
	if err == nil {
		now := time.Now().Format(time.RFC3339)
		conversation.LastMessageAt = &now
		s.conversationRepo.Update(ctx, conversation)
	}

	return toMessageResponse(message), nil, http.StatusCreated
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

	participant, err := s.conversationRepo.GetParticipant(ctx, req.ConversationID, userID)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}

	messages, total, err := s.messageRepo.GetByConversationID(ctx, req.ConversationID, req.Page, req.PageSize, req.Before, participant.JoinedAt)
	if err != nil {
		return nil, 0, err, http.StatusInternalServerError
	}

	responses := make([]model.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = *toMessageResponse(&msg)
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

func (s *MessageService) GetByID(ctx context.Context, id string) (*model.MessageResponse, error) {
	msg, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toMessageResponse(msg), nil
}

func (s *MessageService) ToggleReaction(ctx context.Context, userID, messageID string, req model.ToggleReactionRequest) (*model.ReactionResponse, string, error, int) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, "", errors.New("message not found"), http.StatusNotFound
	}

	// Verify if the user is a participant of the conversation
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, message.ConversationID, userID)
	if err != nil {
		return nil, "", err, http.StatusInternalServerError
	}
	if !isParticipant {
		return nil, "", errors.New("you are not a participant of this conversation"), http.StatusForbidden
	}

	reaction := &entity.MessageReaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     req.Emoji,
	}

	action, err := s.messageRepo.ToggleReaction(ctx, reaction)
	if err != nil {
		return nil, "", err, http.StatusInternalServerError
	}

	var resp *model.ReactionResponse
	if action != "remove" {
		resp = &model.ReactionResponse{
			ID:        reaction.ID,
			MessageID: reaction.MessageID,
			UserID:    reaction.UserID,
			Emoji:     reaction.Emoji,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
	}

	return resp, action, nil, http.StatusOK
}

func toMessageResponse(msg *entity.Message) *model.MessageResponse {
	if msg == nil {
		return nil
	}
	resp := &model.MessageResponse{
		ID:               msg.ID,
		ConversationID:   msg.ConversationID,
		SenderID:         msg.SenderID,
		Type:             string(msg.Type),
		Content:          msg.Content,
		MediaURL:         msg.MediaURL,
		ReplyToMessageID: msg.ReplyToMessageID,
		CreatedAt:        msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        msg.UpdatedAt.Format(time.RFC3339),
	}

	if msg.ReplyToMessage != nil {
		resp.ReplyToMessage = &model.MessageResponse{
			ID:             msg.ReplyToMessage.ID,
			ConversationID: msg.ReplyToMessage.ConversationID,
			SenderID:       msg.ReplyToMessage.SenderID,
			Type:           string(msg.ReplyToMessage.Type),
			Content:        msg.ReplyToMessage.Content,
			MediaURL:       msg.ReplyToMessage.MediaURL,
			CreatedAt:      msg.ReplyToMessage.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      msg.ReplyToMessage.UpdatedAt.Format(time.RFC3339),
		}
	}

	if len(msg.Reactions) > 0 {
		resp.Reactions = make([]model.ReactionResponse, len(msg.Reactions))
		for i, r := range msg.Reactions {
			resp.Reactions[i] = model.ReactionResponse{
				ID:        r.ID,
				MessageID: r.MessageID,
				UserID:    r.UserID,
				Emoji:     r.Emoji,
				CreatedAt: r.CreatedAt.Format(time.RFC3339),
			}
		}
	}

	return resp
}
