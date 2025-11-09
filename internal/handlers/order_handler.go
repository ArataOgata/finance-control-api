package handlers

import (
	"encoding/json"
	respo "go-api/internal/dto/base"
	dto "go-api/internal/dto/order_dto"
	"go-api/internal/service"
	"go-api/internal/validators"
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
		respo.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		log.Info("Error validating create order request", slog.String("error", err.Error()))
		respo.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		order, err := h.service.CreateOrder(tx, &req)
		if err != nil {
			log.Info("Error creating order", slog.String("error", err.Error()))
			return err
		}

		if err := respo.SendJSON(w, order); err != nil {
			log.Info("Error encoding response", slog.String("error", err.Error()))
			respo.SendError(w, http.StatusInternalServerError, err.Error())
		}
		return nil
	})

	if err != nil {
		log.Error("Error creating order", slog.String("error", err.Error()))
		respo.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")

}
