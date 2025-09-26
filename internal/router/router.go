package router

import (
	"go-api/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/ping", handlers.PingHandler)

	return r
}
