package dto

type OrderIDS struct {
	OrderID uint `validate:"gt=0" query:"order_id"`
}
