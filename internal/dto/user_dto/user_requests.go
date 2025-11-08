package userdto

import (
	"strings"
)

type UserRequest struct {
	UserID   uint    `json:"user_id" validate:"required,min=1"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=1,max=25"`
	TgID     *int    `json:"tg_id,omitempty" validate:"omitempty"`
	Balance  *int    `json:"balance,omitempty" validate:"omitempty"`
}

type UserRegisterRequest struct {
	TgID     uint   `json:"tg_id,required" validate:"required,gte=1"`
	Username string `json:"username,required" validate:"required,min=3,max=35"`
	Balance  uint   `json:"balance,omitempty" validate:"omitempty,gte=0"`
}

type UserIDS struct {
	UserID uint `validate:"gt=0" query:"id"`
}

func (r *UserRequest) IsEmpty() bool {
	return r.Username == nil && r.TgID == nil && r.Balance == nil
}

func (r *UserRequest) ToMap(balance *int) map[string]interface{} {
	updates := make(map[string]interface{})

	if r.Username != nil {
		updates["username"] = strings.TrimSpace((*r.Username))
	}

	if r.TgID != nil {
		updates["tg_id"] = *r.TgID
	}

	if r.Balance != nil && balance != nil {
		updates["balance"] = *balance + *r.Balance
	} else {
		updates["balance"] = *r.Balance
	}

	return updates

}
