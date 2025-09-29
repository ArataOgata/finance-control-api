package handlers

import (
	"encoding/json"
	"go-api/internal/service"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateCategory(input.Title, input.Description, uint(input.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idUsStr := r.URL.Query().Get("user_id")
	user_id, _ := strconv.Atoi(idUsStr)

	idCatStr := r.URL.Query().Get("category_id")
	category_id, _ := strconv.Atoi(idCatStr)

	category, err := h.service.GetCategory(uint(category_id), uint(user_id))
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

}
