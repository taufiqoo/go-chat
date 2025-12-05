package main

import (
	"fmt"
	"log"

	"github.com/taufiqoo/go-chat/internal/config"
	"github.com/taufiqoo/go-chat/internal/delivery/http/handler"
	"github.com/taufiqoo/go-chat/internal/delivery/http/router"
	"github.com/taufiqoo/go-chat/internal/delivery/websocket"
	"github.com/taufiqoo/go-chat/internal/repository/repositoryImpl"
	"github.com/taufiqoo/go-chat/internal/service"
	"github.com/taufiqoo/go-chat/pkg/database"
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

	// Setup router
	r := router.SetupRouter(userHandler, messageHandler, wsHandler, &cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
