package repository

import (
	"errors"
	"go-api/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	CreateWithTx(tx *gorm.DB, order *models.Order) error
	GetOrder(tx *gorm.DB, orderID uint, userID uint) (*models.Order, error)
	GetAll(tx *gorm.DB, userId uint) ([]*models.Order, error)
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

func (o *orderRepository) GetAll(tx *gorm.DB, userId uint) ([]*models.Order, error) {
	var orders []*models.Order
	if tx == nil {
		return nil, errors.New("transaction is required")
	}

	err := tx.Where("user_id = ?", userId).Preload("Category").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) GetOrderForCategory(tx *gorm.DB, userId uint, categoryId uint) ([]*models.Order, error) {
	return nil, nil
}

func (o *orderRepository) GetOrder(tx *gorm.DB, orderID uint, userID uint) (*models.Order, error) {
	var order models.Order
	err := tx.Where("order_id = ?", orderID).Where("user_id = ?", userID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}
