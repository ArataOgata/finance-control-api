package validators

import (
	dto "go-api/internal/dto/category_dto"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CategoryValidator struct{}

func (v *CategoryValidator) ValidateUpdateRequest(req *dto.UpdateCategoryRequest) error {
	return validate.Struct(req)
}
