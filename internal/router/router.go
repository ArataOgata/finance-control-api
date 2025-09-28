package router

import (
	"github.com/go-chi/chi/v5"

	"go-api/internal/db"
	"go-api/internal/handlers"
	"go-api/internal/repository"
	"go-api/internal/service"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", handlers.PingHandler)

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// public
	r.Post("/register", userHandler.Register)
	r.Get("/user", userHandler.GetUser)

	return r
}
