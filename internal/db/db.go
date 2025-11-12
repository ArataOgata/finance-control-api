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

	dsn := cfg.Database_URL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Авто-миграция моделей
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate table User: %v", err)
	}

	err = DB.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatalf("Failed to migrate table Category: %v", err)
	}

	err = DB.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalf("Failed to migrate table Order: %v", err)
	}

	fmt.Println("Database migrated")
}
