package handlers

import (
	"encoding/json"
	dto "go-api/internal/dto/order_dto"
	"go-api/internal/service"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type OrderHandler struct {
	db      *gorm.DB
	service service.OrderService
}

func NewOrderHandler(db *gorm.DB, service service.OrderService) *OrderHandler {
	return &OrderHandler{db: db, service: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		resp, err := h.service.CreateOrder(tx, &req)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return nil
	})

	if err != nil {
		log.Printf("error: %v", err)
		h.sendJSONError(w, err.Error(), http.StatusBadRequest)
	}

}

func (h *OrderHandler) sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
