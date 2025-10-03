package models

import "time"

// Order — модель для таблицы orders
type Order struct {
	OrderID     uint      `gorm:"primaryKey;autoIncrement" json:"order_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Description string    `gorm:"size:300" json:"description"`
	Amount      int       `gorm:"not null" json:"amount"`
	Request     []byte    `gorm:"type:jsonb;" json:"-"`
	Response    []byte    `gorm:"type:jsonb;" json:"-"`
	NaviDate    time.Time `gorm:"autoCreateTime" json:"navi_date"`
}
