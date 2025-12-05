package handler

import (
	"net/http"
	"strconv"

	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/service"
	"github.com/taufiqoo/go-chat/internal/utils"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")

	var req domain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message, err := h.messageService.SendMessage(userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Message sent successfully", message)
}

func (h *MessageHandler) GetChatHistory(c *gin.Context) {
	userID := c.GetUint("userID")
	otherUserIDStr := c.Param("userId")

	otherUserID, err := strconv.ParseUint(otherUserIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	messages, err := h.messageService.GetChatHistory(userID, uint(otherUserID), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve chat history")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Chat history retrieved successfully", messages)
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	messageIDStr := c.Param("messageId")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if err := h.messageService.MarkMessageAsRead(uint(messageID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark message as read")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Message marked as read", nil)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.messageService.GetUnreadCount(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Unread count retrieved successfully", gin.H{
		"unread_count": count,
	})
}

func (h *MessageHandler) GetChatList(c *gin.Context) {
	userID := c.GetUint("userID")

	chatList, err := h.messageService.GetChatList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get chat list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": chatList,
	})
}
