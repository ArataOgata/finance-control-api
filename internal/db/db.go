package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-api/config"
	"go-api/internal/models"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) {
	var err error

	dsn := cfg.DatabaseURL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Авто-миграция моделей
	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Order{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migrated")
}
