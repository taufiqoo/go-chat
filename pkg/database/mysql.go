package database

import (
	"fmt"
	"log"

	"github.com/taufiqoo/go-chat/internal/config"
	"github.com/taufiqoo/go-chat/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// Check apakah table sudah ada dari migration files
	if db.Migrator().HasTable(&domain.User{}) && db.Migrator().HasTable(&domain.Message{}) {
		log.Println("✓ Tables already exist, skipping auto migration")
		return nil
	}

	// Kalau belum ada, baru jalankan auto migrate
	log.Println("→ Running auto migration...")
	return db.AutoMigrate(
		&domain.User{},
		&domain.Message{},
	)
}
