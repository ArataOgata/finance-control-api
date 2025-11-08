package handlers

import (
	"encoding/json"
	"go-api/internal/decoder"
	resp "go-api/internal/dto/base_response"
	dto "go-api/internal/dto/category_dto"
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

	category, err := h.service.CreateCategory(req.Title, req.Description, req.UserID)
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

	var ids dto.QueryIDs

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

	category, err := h.service.GetCategory(ids.CategoryID, ids.UserID)
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

	category, err := h.service.UpdateCategory(&req)
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
