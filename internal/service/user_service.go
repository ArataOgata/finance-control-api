package service

import (
	"errors"
	"fmt"
	userdto "go-api/internal/dto/user_dto"
	"go-api/internal/models"
	"go-api/internal/repository"
	"log"

	"gorm.io/gorm"
)

// UserService — интерфейс, определяющий бизнес-логику для работы с пользователями
type UserService interface {
	Register(username string, balance uint, tg_id uint) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	UpdateUser(req *userdto.UserRequest) (*models.User, error)
}

// userService — структура, реализующая интерфейс UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService — конструктор для создания сервиса
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register создаёт нового пользователя с указанным именем и балансом
func (s *userService) Register(username string, balance uint, tg_id uint) (*models.User, error) {

	exists, err := s.repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

		} else {
			log.Println("Failed to check username:", err)
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
	}

	if exists != nil {
		log.Println("Username already taken:", username, exists)
		return nil, errors.New("username already taken") // Ошибка, если имя занято
	}

	user := &models.User{
		Username: username,
		Balance:  int(balance),
		Tg_id:    int(tg_id),
	}

	err = s.repo.Create(user)
	return user, err
}

// GetUser получает пользователя по ID
func (s *userService) GetUser(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(req *userdto.UserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.IsEmpty() {
		return nil, errors.New("no fields to update")
	}

	if *req.Balance < 0 {
		return nil, fmt.Errorf(" <0 : %w", err)
	}
	updates := req.ToMap(&user.Balance)

	if err := s.repo.Update(user, updates); err != nil {
		return nil, fmt.Errorf("failde to update user: %w", err)
	}

	if username, ok := updates["username"].(string); ok {
		user.Username = username
	}

	if tg_id, ok := updates["tg_id"].(int); ok {
		user.Tg_id = tg_id
	}

	if balance, ok := updates["balance"].(int); ok {
		user.Balance = balance
	}

	return user, nil
}
