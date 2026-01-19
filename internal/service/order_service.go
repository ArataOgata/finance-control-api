package service

import (
	"fmt"
	dto "go-api/internal/dto/order_dto"
	"go-api/internal/models"
	"go-api/internal/repository"

	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(userID uint, tx *gorm.DB, req *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetAllGroupedByCategory(tx *gorm.DB, userId uint) (dto.OrdersByCategoryDTO, error)
	GetOrderByid(tx *gorm.DB, orderID uint, userID uint) (*dto.OrderResponse, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
	categRepo repository.CategoryRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	categRepo repository.CategoryRepository,
) OrderService {
	return &orderService{orderRepo: orderRepo, userRepo: userRepo, categRepo: categRepo}
}

func (os *orderService) CreateOrder(userID uint, tx *gorm.DB, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {

	user, err := os.userRepo.FindByIDWithTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	category, err := os.categRepo.GetByCategoryIdWithTx(tx, req.CategoryID, user.UserID) // Получение категории с tx
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if user.Balance < req.Amount {
		return nil, fmt.Errorf("failed, insufficient funds")
	}

	user.Balance -= req.Amount
	if err := os.userRepo.UpdateWithTx(tx, user); err != nil {
		return nil, fmt.Errorf("недостаточно средств. Доступно: %d, требуется: %d", user.Balance, req.Amount)
	}

	order := &models.Order{ // Создание модели заказа из DTO
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
	}

	if err := os.orderRepo.CreateWithTx(tx, order); err != nil { // Создание заказа с tx
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	category.Total += req.Amount                                    // Обновление Total в памяти
	if err := os.categRepo.UpdateWithTx(tx, category); err != nil { // Обновление категории с tx
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &dto.OrderResponse{
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		CategoryID:  order.CategoryID,
		Description: order.Description,
		Amount:      order.Amount,
		NaviDate:    order.NaviDate,
	}, nil

}

func (os *orderService) GetAllGroupedByCategory(tx *gorm.DB, userId uint) (dto.OrdersByCategoryDTO, error) {
	user, err := os.userRepo.FindByIDWithTx(tx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	orders, err := os.orderRepo.GetAll(tx, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	grouped := make(map[string][]*models.Order)
	for _, order := range orders {
		grouped[order.Category.Title] = append(grouped[order.Category.Title], order)
	}

	result := make(dto.OrdersByCategoryDTO)
	for title, orders := range grouped {
		items := make([]*dto.OrderItemDTO, len(orders))
		for i, o := range orders {
			items[i] = &dto.OrderItemDTO{
				OrderID:     o.OrderID,
				Amount:      o.Amount,
				Description: o.Description,
				NaviDate:    o.NaviDate,
			}
		}
		result[title] = items
	}

	return result, nil
}

func (os *orderService) GetOrderByid(tx *gorm.DB, orderID uint, userID uint) (*dto.OrderResponse, error) {
	user, err := os.userRepo.FindByIDWithTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	order, err := os.orderRepo.GetOrder(tx, orderID, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &dto.OrderResponse{
		OrderID:     order.OrderID,
		UserID:      order.UserID,
		CategoryID:  order.CategoryID,
		Description: order.Description,
		Amount:      order.Amount,
		NaviDate:    order.NaviDate,
	}, nil

}
