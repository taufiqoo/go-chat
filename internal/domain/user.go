package domain

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Fullname  string    `json:"fullname" gorm:"type:varchar(50);not null"`
	Photo     string    `json:"photo" gorm:"type:varchar(255)"`
	Username  string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	Email     string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Fullname string `json:"fullname" binding:"required,min=3,max=100"`
	Photo    string `json:"photo"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Photo    string `json:"photo"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
}
