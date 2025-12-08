package router

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/taufiqoo/go-chat/internal/config"
	"github.com/taufiqoo/go-chat/internal/delivery/http/handler"
	"github.com/taufiqoo/go-chat/internal/delivery/http/middleware"
	"github.com/taufiqoo/go-chat/internal/delivery/websocket"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
}

func SetupRouter(
	userHandler *handler.UserHandler,
	messageHandler *handler.MessageHandler,
	wsHandler *websocket.Handler,
	cfg *config.Config,
	redisClient *redis.Client,
	rateLimitConfig RateLimitConfig,
) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	if redisClient != nil {
		rateLimiter := middleware.NewRateLimiter(
			redisClient,
			rateLimitConfig.MaxRequests,
			rateLimitConfig.Window,
		)
		r.Use(rateLimiter.Middleware())
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User routes
			protected.GET("/profile", userHandler.GetProfile)
			protected.GET("/users", userHandler.GetAllUsers)

			// Chat routes
			protected.POST("/messages", messageHandler.SendMessage)
			protected.GET("/messages/:userId", messageHandler.GetChatHistory)
			protected.PATCH("/messages/:messageId/read", messageHandler.MarkAsRead)
			protected.GET("/messages/unread/count", messageHandler.GetUnreadCount)
			protected.GET("/messages/chat-list", messageHandler.GetChatList)
		}

		// WebSocket route
		api.GET("/ws", wsHandler.HandleWebSocket)
	}

	return r
}
