package dto

import (
	"strings"
)

type UpdateCategoryRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required,min=1"`
	UserID      uint    `json:"user_id" validate:"required,min=1"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=1000"`
	Total       *int    `json:"total,omitempty" validate:"omitempty,min=0"`
}

type CreateCategory struct {
	Title       string `json:"title" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"required,min=1,max=300"`
	UserID      uint   `json:"user_id" validate:"required,gt=1"`
}

func (r *UpdateCategoryRequest) IsEmpty() bool {
	return r.Title == nil && r.Description == nil && r.Total == nil
}

func (r *UpdateCategoryRequest) ToMap() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.Title != nil {
		updates["title"] = strings.TrimSpace(*r.Title)
	}
	if r.Description != nil {
		updates["description"] = strings.TrimSpace(*r.Description)
	}
	if r.Total != nil {
		updates["total"] = *r.Total
	}
	return updates
}
