package router

import (
	"go-api/internal/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-api/internal/db"
	"go-api/internal/handlers"
	"go-api/internal/repository"
	"go-api/internal/service"
)

func NewRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	catRepo := repository.NewCategoryRepository(db.DB)
	catService := service.NewCategoryService(catRepo, userRepo)
	catHandler := handlers.NewCategoryHandler(catService)

	ordRepo := repository.NewOrderRepository(db.DB)
	ordService := service.NewOrederService(ordRepo, userRepo, catRepo)
	ordHandler := handlers.NewOrderHandler(db.DB, ordService)

	r.Use(middleware.RequestID)
	r.Use(logger.MiddleLogger(log))

	r.Get("/ping", handlers.PingHandler)

	// public
	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Get("/", userHandler.GetUser)
		r.Patch("/update", userHandler.UpdateUser)
	})

	r.Route("/api/v1/category", func(r chi.Router) {
		r.Post("/", catHandler.CreateCategory)
		r.Get("/", catHandler.GetCategory)
		r.Patch("/update", catHandler.UpdateCategory)
	})

	r.Route("/api/v1/order", func(r chi.Router) {
		r.Post("/", ordHandler.CreateOrder)
	})

	return r
}
