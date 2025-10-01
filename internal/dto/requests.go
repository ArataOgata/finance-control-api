package dto

import (
	"errors"
	"strings"
)

type UpdateCategoryRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required"`
	UserID      uint    `json:"user_id" validate:"required"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Total       *int    `json:"total,omitempty" validate:"omitempty,min=0"`
}

func (r *UpdateCategoryRequest) Validate() error {
	if r.Total != nil && *r.Total < 0 {
		return errors.New("total cannot be negative")
	}
	return nil
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
