package userdto

import "strings"

type UpdateUserRequest struct {
	UserID   uint    `json:"user_id" validate:"required,min=1"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=1,max=25"`
	Tg_id    *int    `json:"tg_id,omitempty" validate:"omitempty"`
	Balance  *int    `json:"balance,omitempty" validate:"omitempty,min=0"`
}

func (r *UpdateUserRequest) IsEmpty() bool {
	return r.Username == nil && r.Tg_id == nil && r.Balance == nil
}

func (r *UpdateUserRequest) ToMap() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.Username != nil {
		updates["username"] = strings.TrimSpace((*r.Username))
	}
	if r.Tg_id != nil {
		updates["tg_id"] = *r.Tg_id
	}
	if r.Balance != nil {
		updates["balance"] = *r.Balance
	}

	return updates
}
