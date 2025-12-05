package repository

import "github.com/taufiqoo/go-chat/internal/domain"

type MessageRepository interface {
	Create(message *domain.Message) error
	FindByID(id uint) (*domain.Message, error)
	GetChatHistory(userID, otherUserID uint, limit int) ([]domain.Message, error)
	MarkAsRead(messageID uint) error
	GetUnreadCount(userID uint) (int64, error)
	GetChatList(userID uint) ([]domain.ChatListItem, error)
	GetUnreadCountByUser(currentUserID, otherUserID uint) (int64, error)
}
