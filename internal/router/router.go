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

	catRepo := repository.NewCategoryRepository(db.DB)
	catService := service.NewCategoryService(catRepo, userRepo)
	catHandler := handlers.NewCategoryHandler(catService)

	// public
	r.Post("/register", userHandler.Register)
	r.Get("/user", userHandler.GetUser)

	r.Post("/category", catHandler.CreateCategory)
	r.Get("/category", catHandler.GetCategory)
	r.Patch("/updateCat", catHandler.UpdateCategory)

	return r
}
