package models

import "time"

// Category — модель для таблицы categories
type Category struct {
	CategoryID  uint      `gorm:"primaryKey;autoIncrement" json:"category_id"`
	Title       string    `gorm:"size:60;unique;not null" json:"title"`
	Description string    `gorm:"size:300" json:"description"`
	Total       int       `gorm:"default:0" json:"total"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UserID      uint      `gorm:"not null" json:"user_id"`
}
