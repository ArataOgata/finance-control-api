package validators

import (
	"errors"
	"go-api/internal/dto"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CategoryValidator struct{}

func (v *CategoryValidator) ValidateIDs(categoryID, userID uint) error {
	if userID == 0 {
		return errors.New("user ID is required")
	}
	if categoryID == 0 {
		return errors.New("category ID is required")
	}
	return nil
}

func (v *CategoryValidator) ValidateUpdateRequest(req *dto.UpdateCategoryRequest) error {
	return validate.Struct(req)
}
