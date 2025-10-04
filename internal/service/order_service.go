package service

import (
	"fmt"
	dto "go-api/internal/dto/order_dto"
	"go-api/internal/models"
	"go-api/internal/repository"

	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(tx *gorm.DB, req *dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
	categRepo repository.CategoryRepository
}

func NewOrederService(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	categRepo repository.CategoryRepository,
) OrderService {
	return &orderService{orderRepo: orderRepo, userRepo: userRepo, categRepo: categRepo}
}

func (os *orderService) CreateOrder(tx *gorm.DB, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// log.Info().Uint("user_id", req.UserID).Uint("category_id", req.CategoryID).Msg("Creating order")

	user, err := os.userRepo.FindByIDWithTx(tx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	category, err := os.categRepo.GetByCategoryIdWithTx(tx, req.CategoryID, user.UserID) // Получение категории с tx
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	order := &models.Order{ // Создание модели заказа из DTO
		UserID:      req.UserID,
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
