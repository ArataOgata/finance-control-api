package models

import "time"

// Order — модель для таблицы orders
type Order struct {
	OrderID     uint      `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Description string    `gorm:"size:300" json:"description"`
	Request     []byte    `gorm:"type:jsonb;not null" json:"request"`
	Response    []byte    `gorm:"type:jsonb;not null" json:"response"`
	NaviDate    time.Time `gorm:"autoCreateTime" json:"navi_date"`
}
