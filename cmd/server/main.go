package main

import (
	"fmt"
	"log"
	"time"

	"github.com/taufiqoo/go-chat/internal/config"
	"github.com/taufiqoo/go-chat/internal/delivery/http/handler"
	"github.com/taufiqoo/go-chat/internal/delivery/http/router"
	"github.com/taufiqoo/go-chat/internal/delivery/websocket"
	"github.com/taufiqoo/go-chat/internal/repository/repositoryImpl"
	"github.com/taufiqoo/go-chat/internal/service"
	"github.com/taufiqoo/go-chat/pkg/database"
	redisClient "github.com/taufiqoo/go-chat/pkg/redis"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewMySQLConnection(&cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	redis := redisClient.NewRedisClient(&cfg)
	if redis == nil {
		log.Println(" Redis not available - rate limiting disabled")
	}

	// Initialize repositories
	userRepo := repositoryImpl.NewUserRepository(db)
	messageRepo := repositoryImpl.NewMessageRepository(db)

	// Initialize usecases
	userService := service.NewUserService(userRepo, &cfg)
	messageService := service.NewMessageService(messageRepo, userRepo)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	userHandler := handler.NewsUserHandler(userService)
	messageHandler := handler.NewMessageHandler(messageService)
	wsHandler := websocket.NewHandler(hub, messageService)

	// Rate limiter config
	rateLimitConfig := router.RateLimitConfig{
		MaxRequests: 1,
		Window:      2 * time.Minute,
	}

	// Setup router
	r := router.SetupRouter(
		userHandler,
		messageHandler,
		wsHandler,
		&cfg,
		redis,
		rateLimitConfig,
	)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
