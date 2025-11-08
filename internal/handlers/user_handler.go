package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	resp "go-api/internal/dto/base_response"
	userdto "go-api/internal/dto/user_dto"
	"go-api/internal/service"
	"go-api/internal/validators"
	//_ "go-api/internal/validators"

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
		Balance  uint   `json:"balance"`
		Tg_id    uint   `json:"tgID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	//if err := validators.ValidateUsername(input.Username); err != nil {
	//	http.Error(w, err.validators(), http.StatusBadRequest)
	//	return
	//}

	user, err := h.service.Register(input.Username, input.Balance, input.Tg_id)
	if err != nil {
		log.Info("Error registering user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       user,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("Successfully registered user", slog.String("username", user.Username))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.GetUser"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn("Invalid user ID provided", slog.String("id", idStr), slog.String("error", "Invalid user ID"))
		resp.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		log.Info("Error getting user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       user,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("Successfully fetched user", slog.String("username", user.Username))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.UpdateUser"
	var req userdto.UserRequest

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	validator := &validators.UserValidator{}

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

	user, err := h.service.UpdateUser(&req)
	if err != nil {
		log.Info("Error updating user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.Response{
		StatusCode: resp.StatusOK,
		Data:       user,
	}); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Successfully updated user", slog.String("username", user.Username))

}
