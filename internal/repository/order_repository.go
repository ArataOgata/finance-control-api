package repository

import (
	"errors"
	"go-api/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	CreateWithTx(tx *gorm.DB, order *models.Order) error
	FindByID(tx *gorm.DB, userId uint, categoryId uint, orderId uint) (*models.Order, error)
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

func (o *orderRepository) FindByID(tx *gorm.DB, userId uint, categoryId uint, orderId uint) (*models.Order, error) {
	var order models.Order
	if tx == nil {
		return nil, errors.New("transaction is required")
	}
	err := o.db.Where("user_id = ? AND category_id = ? AND order_id = ?", userId, categoryId, orderId).First(&order).Error
	return &order, err
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
