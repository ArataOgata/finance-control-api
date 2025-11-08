package dto

type CategoryResponse struct {
	CategoryID  uint    `json:"category_id" validate:"required,min=1"`
	UserID      uint    `json:"user_id" validate:"required,min=1"`
	Title       string  `json:"title,omitempty" validate:"min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=1000"`
	Total       int     `json:"total,omitempty" validate:"min=0"`
}
