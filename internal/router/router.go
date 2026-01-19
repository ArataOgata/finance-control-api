package router

import (
	midlle "go-api/internal/middleware"
	"go-api/internal/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-api/internal/db"
	"go-api/internal/handlers"
	"go-api/internal/repository"
	"go-api/internal/service"
)

func NewRouter(log *slog.Logger, secretKey string) *chi.Mux {
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, log)

	catRepo := repository.NewCategoryRepository(db.DB)
	catService := service.NewCategoryService(catRepo, userRepo)
	catHandler := handlers.NewCategoryHandler(catService, log)

	ordRepo := repository.NewOrderRepository(db.DB)
	ordService := service.NewOrderService(ordRepo, userRepo, catRepo)
	ordHandler := handlers.NewOrderHandler(db.DB, ordService, log)

	r.Use(middleware.RequestID)
	r.Use(logger.MiddleLogger(log))

	r.Get("/ping", handlers.PingHandler)

	// public

	r.Group(func(r chi.Router) {
		r.Route("/api/v1/", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(midlle.AuthMiddleware(secretKey))

		r.Route("/api/v1/user", func(r chi.Router) {
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
			r.Get("/all", ordHandler.GetAllOrders)
			r.Get("/", ordHandler.GetOrderByID)
		})
	})

	return r
}
