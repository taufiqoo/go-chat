package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/service"
	"github.com/taufiqoo/go-chat/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub            *Hub
	messageService service.MessageService
}

func NewHandler(hub *Hub, messageService service.MessageService) *Handler {
	return &Handler{
		hub:            hub,
		messageService: messageService,
	}
}

type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:      h.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   uint(userID),
		messages: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go h.handleMessages(client)
	client.readPump()
}

// func (h *Handler) handleMessages(client *Client) {
// 	for {
// 		_, msg, err := client.conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading WS message:", err)
// 			client.hub.unregister <- client
// 			client.conn.Close()
// 			return
// 		}

// 		var wsMsg WSMessage
// 		if err := json.Unmarshal(msg, &wsMsg); err != nil {
// 			log.Println("Invalid WS message:", err)
// 			continue
// 		}

// 		// Simpan pesan lewat service
// 		savedMsg, err := h.messageService.SendMessage(client.userID, &domain.SendMessageRequest{
// 			ReceiverID: wsMsg.ReceiverID,
// 			Content:    wsMsg.Content,
// 		})
// 		if err != nil {
// 			log.Println("Error saving message:", err)
// 			continue
// 		}

// 		// Broadcast ke receiver
// 		messageBytes, _ := json.Marshal(savedMsg)
// 		h.hub.BroadcastToUser(wsMsg.ReceiverID, messageBytes)
// 	}
// }

func (h *Handler) handleMessages(client *Client) {
	for msg := range client.messages { // Baca dari channel, bukan dari conn
		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Println("Invalid WS message:", err)
			continue
		}

		// Simpan pesan ke database
		savedMsg, err := h.messageService.SendMessage(client.userID, &domain.SendMessageRequest{
			ReceiverID: wsMsg.ReceiverID,
			Content:    wsMsg.Content,
		})
		if err != nil {
			log.Println("Error saving message:", err)
			continue
		}

		// Broadcast ke receiver
		messageBytes, _ := json.Marshal(savedMsg)
		h.hub.BroadcastToUser(wsMsg.ReceiverID, messageBytes)

		// Broadcast juga ke sender (echo back)
		h.hub.BroadcastToUser(client.userID, messageBytes)
	}
}
