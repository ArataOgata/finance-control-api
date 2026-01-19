package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	resp "go-api/internal/dto/base"
	userdto "go-api/internal/dto/user_dto"
	midlle "go-api/internal/middleware"
	"go-api/internal/service"
	"go-api/internal/validators"

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

	var req userdto.UserRegisterRequest

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("Error decoding body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		log.Info("Error validating body", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(req.Username, req.Balance, req.TgID)
	if err != nil {
		log.Info("Error registering user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := resp.SendJSON(w, user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
	}
	log.Info("Successfully registered user", slog.String("username", user.Username))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "Handler/UserHandler.GetUser"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := midlle.UserIDFromRequest(r)
	if !ok {
		resp.SendError(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		log.Info("Error getting user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := resp.SendJSON(w, user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
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

	if err := validators.Validate.Struct(&req); err != nil {
		log.Info("Error validating update request", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		log.Info("Error updating user", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := resp.SendJSON(w, user); err != nil {
		log.Info("Error encoding response", slog.String("error", err.Error()))
		resp.SendError(w, http.StatusInternalServerError, err.Error())
	}

	log.Info("Successfully updated user", slog.String("username", user.Username))

}
