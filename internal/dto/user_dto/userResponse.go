package userdto

type UserResponse struct {
	Username string `json:"username"`
	TgID     uint   `json:"tg_id"`
	Balance  int    `json:"balance"`
}
