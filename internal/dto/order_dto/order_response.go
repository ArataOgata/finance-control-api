package dto

import "time"

type OrderResponse struct {
	OrderID     uint      `json:"order_id"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	NaviDate    time.Time `json:"navi_date"`
}
