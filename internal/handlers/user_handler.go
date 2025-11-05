package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	userdto "go-api/internal/dto/user_dto"
	"go-api/internal/service"
	"go-api/internal/validators"
	// _ "go-api/internal/validation"

	"github.com/go-chi/chi/v5/middleware"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service service.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  log,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.Register"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var input struct {
		Username string `json:"username"`
		Balance  int    `json:"balance"`
		Tg_id    int    `json:"tgID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// if err := validation.ValidateUsername(input.Username); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	user, err := h.service.Register(input.Username, input.Balance, input.Tg_id)
	if err != nil {
		log.Info("Error registering user", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info("Successfully registered user", slog.String("username", user.Username))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.GetUser"
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		log.Info("Error getting user", slog.String("error", err.Error()))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info("Successfully fetched user", slog.String("username", user.Username))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.UpdateUser"
	var req userdto.UpdateUserRequest

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	validator := &validators.UserValidator{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateUpdateRequest(&req); err != nil {
		log.Info("Error validating update request", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(&req)
	if err != nil {
		log.Info("Error updating user", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Successfully updated user", slog.String("username", user.Username))

}
