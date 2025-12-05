package repositoryImpl

import (
	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/repository"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) FindByID(id uint) (*domain.Message, error) {
	var message domain.Message
	err := r.db.Preload("Sender").Preload("Receiver").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetChatHistory(userID, otherUserID uint, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	err := r.db.
		Preload("Sender").
		Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, err
}

func (r *messageRepository) GetChatList(userID uint) ([]domain.ChatListItem, error) {
	var chatList []domain.ChatListItem

	err := r.db.Raw(`
		SELECT 
			u.id AS user_id,
			u.fullname,
			u.photo,
			u.username,
			m.content,
			m.is_read,
			m.created_at,
			CASE WHEN m.sender_id = ? THEN true ELSE false END AS is_sender
		FROM users u
		JOIN (
			SELECT 
				CASE 
					WHEN sender_id = ? THEN receiver_id 
					ELSE sender_id 
				END AS contact_id,
				MAX(created_at) AS last_message_time
			FROM messages
			WHERE sender_id = ? OR receiver_id = ?
			GROUP BY contact_id
		) latest ON latest.contact_id = u.id
		JOIN messages m ON (
			(m.sender_id = ? AND m.receiver_id = u.id) OR 
			(m.sender_id = u.id AND m.receiver_id = ?)
		) AND m.created_at = latest.last_message_time
		ORDER BY m.created_at DESC
	`, userID, userID, userID, userID, userID, userID).Scan(&chatList).Error

	if err != nil {
		return nil, err
	}

	for i := range chatList {
		count, _ := r.GetUnreadCountByUser(userID, chatList[i].UserID)
		chatList[i].UnreadCount = int(count)
	}

	return chatList, nil
}

func (r *messageRepository) MarkAsRead(messageID uint) error {
	return r.db.Model(&domain.Message{}).Where("id = ?", messageID).Update("is_read", true).Error
}

func (r *messageRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *messageRepository) GetUnreadCountByUser(currentUserID, otherUserID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?",
			otherUserID, currentUserID, false).
		Count(&count).Error
	return count, err
}
