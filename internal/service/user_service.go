package service // Объявляем пакет service, который находится в internal/service

import (
	"errors" // Импортируем стандартный пакет для создания ошибок
	"fmt"
	userdto "go-api/internal/dto/user_dto"
	"go-api/internal/models"     // Импортируем пакет с моделью User
	"go-api/internal/repository" // Импортируем пакет repository для работы с UserRepository
	"log"

	"gorm.io/gorm"
)

// UserService — интерфейс, определяющий бизнес-логику для работы с пользователями
type UserService interface {
	Register(username string, balance int, tg_id int) (*models.User, error) // Регистрация нового пользователя
	GetUser(id uint) (*models.User, error)                                  // Получение данных пользователя по ID
	UpdateUser(req *userdto.UpdateUserRequest) (*models.User, error)
}

// userService — структура, реализующая интерфейс UserService
type userService struct {
	repo repository.UserRepository // Зависимость от репозитория для доступа к данным
}

// NewUserService — конструктор для создания сервиса
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo} // Создаём структуру userService с переданным репозиторием
}

// Register создаёт нового пользователя с указанным именем и балансом
func (s *userService) Register(username string, balance int, tg_id int) (*models.User, error) {
	// Проверяем, не занято ли имя пользователя
	exists, err := s.repo.FindByUsername(username) // Игнорируем ошибку (неидеально)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь не найден, можно продолжать регистрацию
		} else {
			log.Println("Failed to check username:", err)               // Логируем ошибку
			return nil, fmt.Errorf("failed to check username: %w", err) // Возвращаем обёрнутую ошибку
		}
	}

	if exists != nil {
		log.Println("Username already taken:", username, exists)
		return nil, errors.New("username already taken") // Ошибка, если имя занято
	}

	user := &models.User{
		Username: username, // Устанавливаем имя пользователя
		Balance:  balance,  // Устанавливаем начальный баланс
		Tg_id:    tg_id,
	}

	err = s.repo.Create(user) // Сохраняем пользователя в базе через репозиторий
	return user, err          // Возвращаем созданного пользователя или ошибку
}

// GetUser получает пользователя по ID
func (s *userService) GetUser(id uint) (*models.User, error) {
	return s.repo.FindByID(id) // Вызываем метод репозитория для получения пользователя
}

func (s *userService) UpdateUser(req *userdto.UpdateUserRequest) (*models.User, error) {
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
