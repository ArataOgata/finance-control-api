package validators

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type CategoryValidator struct{}

type UserValidator struct{}

type IDs struct{}
