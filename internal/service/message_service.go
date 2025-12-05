package service

import (
	"errors"

	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/repository"
)

type MessageService interface {
	SendMessage(senderID uint, req *domain.SendMessageRequest) (*domain.Message, error)
	GetChatHistory(userID, otherUserID uint, limit int) ([]domain.Message, error)
	MarkMessageAsRead(messageID uint) error
	GetUnreadCount(userID uint) (int64, error)
	GetChatList(userID uint) ([]domain.ChatListItem, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
}

func NewMessageService(messageRepo repository.MessageRepository, userRepo repository.UserRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (c *messageService) SendMessage(senderID uint, req *domain.SendMessageRequest) (*domain.Message, error) {
	_, err := c.userRepo.FindByID(req.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}

	message := &domain.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	if err := c.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return c.messageRepo.FindByID(message.ID)
}

func (c *messageService) GetChatHistory(userID, otherUserID uint, limit int) ([]domain.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return c.messageRepo.GetChatHistory(userID, otherUserID, limit)
}

func (c *messageService) MarkMessageAsRead(messageID uint) error {
	return c.messageRepo.MarkAsRead(messageID)
}

func (c *messageService) GetUnreadCount(userID uint) (int64, error) {
	return c.messageRepo.GetUnreadCount(userID)
}

func (s *messageService) GetChatList(userID uint) ([]domain.ChatListItem, error) {
	return s.messageRepo.GetChatList(userID)
}
