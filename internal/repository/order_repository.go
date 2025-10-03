package repository

import (
	"errors"
	"go-api/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	CreateWithTx(tx *gorm.DB, order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) Create(order *models.Order) error {
	return o.db.Create(order).Error
}

func (o *orderRepository) CreateWithTx(tx *gorm.DB, order *models.Order) error {
	if tx == nil {
		return errors.New("transaction is required")
	}

	return tx.Create(order).Error
}
