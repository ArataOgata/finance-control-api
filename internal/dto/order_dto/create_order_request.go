package dto

type CreateOrderRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	CategoryID  uint   `json:"category_id" validate:"required"`
	Description string `json:"description" validate:"required,max=300"`
	Amount      int    `json:"amount" validate:"required,min=1"`
}
