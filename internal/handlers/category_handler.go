package handlers

import (
	"encoding/json"
	"go-api/internal/decoder"
	resp "go-api/internal/dto/base"
	dto "go-api/internal/dto/category_dto"
	midlle "go-api/internal/middleware"
	"go-api/internal/service"
	"go-api/internal/validators"
	"log/slog"
	"net/http"

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
	const op = "Handler/CategoryHandler.CreateCategory"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.CreateCategory
	userID, ok := midlle.UserIDFromRequest(r)
	if !ok {
		resp.SendError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		log.Info("Error validating ", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.CreateCategory(req.Title, req.Description, userID)
	if err != nil {
		log.Info("Error creating category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := resp.SendJSON(w, category); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/CategoryHandler.GetCategory"
	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var ids resp.QueryIDs
	userID, ok := midlle.UserIDFromRequest(r)
	if !ok {
		resp.SendError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	if err := decoder.Decoder.Decode(&ids, r.URL.Query()); err != nil {
		log.Info("Invalid query parameters", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if err := validators.Validate.Struct(ids); err != nil {
		log.Info("Error validating ids", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.GetCategory(ids.CategoryID, userID)
	if err != nil {
		log.Info("Error getting category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := resp.SendJSON(w, category); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("Request completed successfully")

}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/CategoryHandler.UpdateCategory"
	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.UpdateCategoryRequest

	userID, ok := midlle.UserIDFromRequest(r)
	if !ok {
		resp.SendError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		log.Info("Error validating update request", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.UpdateCategory(userID, &req)
	if err != nil {
		log.Info("Error updating category", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := resp.SendJSON(w, category); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Request completed successfully")
}
