package domain

import (
	"time"
)

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`

	Sender   User `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver User `json:"receiver" gorm:"foreignKey:ReceiverID"`
}

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,min=1"`
}

type MessageResponse struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	Sender     UserInfo  `json:"sender"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type ChatListItem struct {
	UserID      uint      `json:"user_id"`
	Fullname    string    `json:"fullname"`
	Photo       string    `json:"photo"`
	Username    string    `json:"username"`
	Content     string    `json:"content"`
	IsRead      bool      `json:"is_read"`
	UnreadCount int       `json:"unread_count"`
	CreatedAt   time.Time `json:"created_at"`
	IsSender    bool      `json:"is_sender"`
}
