package main

import (
	"go-api/internal/logger"
	"log"
	"log/slog"
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
	logs := logger.SetupLogger(cfg.Env)

	db.ConnectDatabase(cfg)

	r := router.NewRouter(logs)

	logs.Info("Server running on ", slog.String("Address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logs.Error("Error starting server: %v", err)
	}

	logs.Info("Server stopped")
}
