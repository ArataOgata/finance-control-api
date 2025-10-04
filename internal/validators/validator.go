package validators

import (
	dto "go-api/internal/dto/category_dto"
	userdto "go-api/internal/dto/user_dto"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CategoryValidator struct{}

type UserValidator struct{}

func (v *CategoryValidator) ValidateUpdateRequest(req *dto.UpdateCategoryRequest) error {
	return validate.Struct(req)
}

func (v *UserValidator) ValidateUpdateRequest(req *userdto.UpdateUserRequest) error {
	return validate.Struct(req)
}
