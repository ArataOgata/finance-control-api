package service // Объявляем пакет service, который находится в internal/service

import (
	"errors" // Импортируем стандартный пакет для создания ошибок
	"fmt"
	"go-api/internal/models"     // Импортируем пакет с моделью User
	"go-api/internal/repository" // Импортируем пакет repository для работы с UserRepository
	"log"

	"gorm.io/gorm"
)

// UserService — интерфейс, определяющий бизнес-логику для работы с пользователями
type UserService interface {
	Register(username string, balance int, tg_id int) (*models.User, error) // Регистрация нового пользователя
	GetUser(id uint) (*models.User, error)                                  // Получение данных пользователя по ID
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
