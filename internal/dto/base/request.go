package base

type QueryIDs struct {
	UserID     uint `validate:"gt=0" query:"user_id"`
	CategoryID uint `validate:"gt=0" query:"category_id"`
}
