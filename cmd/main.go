package main

import (
	"log"
	"net/http"

	"go-api/config"
	"go-api/internal/db"
	"go-api/internal/router"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db.ConnectDatabase(cfg)

	r := router.NewRouter()

	log.Printf("Server running on %s\n", cfg.ServerPort)

	err = http.ListenAndServe(cfg.ServerPort, r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
