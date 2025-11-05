package handlers

import (
	"encoding/json"
	dto "go-api/internal/dto/order_dto"
	"go-api/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db      *gorm.DB
	service service.OrderService
	logger  *slog.Logger
}

func NewOrderHandler(db *gorm.DB, service service.OrderService, log *slog.Logger) *OrderHandler {
	return &OrderHandler{db: db, service: service, logger: log}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.OrderHandler.CreateOrder"
	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding create order request", slog.String("error", err.Error()))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		resp, err := h.service.CreateOrder(tx, &req)
		if err != nil {
			log.Info("Error creating order", slog.String("error", err.Error()))
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Info("Error encoding create order response", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	if err != nil {
		log.Error("Error creating order", slog.String("error", err.Error()))
		h.sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Request completed successfully")

}

func (h *OrderHandler) sendJSONError(w http.ResponseWriter, message string, status int) {
	log := h.logger.With(slog.String("op", "sendJSONError"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Error("Error encoding JSON response", slog.String("error", err.Error()))
	}
}
