package handlers

import (
	"encoding/json"
	resp "go-api/internal/dto/base_response"
	dto "go-api/internal/dto/category_dto"
	"go-api/internal/service"
	"go-api/internal/validators"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

type CategoryHandler struct {
	service service.CategoryService
	logger  *slog.Logger
}

func NewCategoryHandler(service service.CategoryService, log *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		service: service,
		logger:  log,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.CreateCategory"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.CreateCategory(input.Title, input.Description, uint(input.UserID))
	if err != nil {
		log.Info("Error creating category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       category,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.GetCategory"
	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idUsStr := r.URL.Query().Get("user_id")
	user_id, _ := strconv.Atoi(idUsStr)

	idCatStr := r.URL.Query().Get("category_id")
	category_id, _ := strconv.Atoi(idCatStr)

	category, err := h.service.GetCategory(uint(category_id), uint(user_id))
	if err != nil {
		log.Info("Error getting category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       category,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")

}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.UpdateCategory"
	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.UpdateCategoryRequest

	validator := &validators.CategoryValidator{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.ValidateUpdateRequest(&req); err != nil {
		log.Info("Error validating update request", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.UpdateCategory(&req)
	if err != nil {
		log.Info("Error updating category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       category,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")
}
