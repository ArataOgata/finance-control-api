package repository // Объявляем пакет repository, который находится в директории internal/repository

import (
	"errors"
	"go-api/internal/models" // Импортируем пакет с моделью User, которая описывает структуру данных пользователя

	"gorm.io/gorm" // Импортируем GORM — библиотеку для работы с базой данных
)

// UserRepository — интерфейс, определяющий контракт для работы с пользователями в базе данных
type UserRepository interface {
	Create(user *models.User) error                       // Метод для создания нового пользователя
	FindByID(id uint) (*models.User, error)               // Метод для поиска пользователя по ID
	FindByUsername(username string) (*models.User, error) // Метод для поиска пользователя по имени
	Update(user *models.User) error                       // Метод для обновления данных пользователя
	Delete(id uint) error                                 // Метод для удаления пользователя по ID
}

// userRepository — структура, которая реализует интерфейс UserRepository
type userRepository struct {
	db *gorm.DB // Поле для хранения подключения к базе данных (GORM)
}

// NewUserRepository — конструктор для создания экземпляра репозитория
// Принимает подключение к базе данных и возвращает интерфейс UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db} // Создаём структуру userRepository с переданным подключением и возвращаем её как интерфейс
}

// Create сохраняет нового пользователя в базе данных
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error // Используем GORM для создания записи в базе данных
}

// FindByID ищет пользователя по его ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User               // Создаём переменную для хранения результата запроса
	err := r.db.First(&user, id).Error // Ищем первую запись в базе с указанным ID
	return &user, err                  // Возвращаем указатель на пользователя и ошибку (если есть)
}

// FindByUsername ищет пользователя по имени
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Возвращаем nil вместо &user
		}
		return nil, err // Возвращаем nil для других ошибок
	}
	return &user, nil // Возвращаем указатель на пользователя, если он найден                                             // Возвращаем указатель на пользователя и ошибку
}

// Update обновляет данные пользователя в базе
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error // Используем GORM для обновления записи
}

// Delete удаляет пользователя по ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error // Удаляем запись из базы по ID
}
