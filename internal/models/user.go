package models

type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username string `gorm:"size:25;unique;not null" json:"username"`
	Balance  int    `gorm:"default:0" json:"balance"`
}
