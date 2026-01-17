package dto

type OrderIDS struct {
	UserID  uint `validate:"gt=0" query:"user_id"`
	OrderID uint `validate:"gt=0" query:"order_id"`
}
