package base

type QueryIDs struct {
	CategoryID uint `validate:"gt=0" query:"category_id"`
}
