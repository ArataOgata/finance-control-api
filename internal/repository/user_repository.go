package repository

import (
	"errors"
	"go-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByIDWithTx(tx *gorm.DB, id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Update(user *models.User, updates map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB // Поле для хранения подключения к базе данных (GORM)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db} // Создаём структуру userRepository с переданным подключением и возвращаем её как интерфейс
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Возвращаем nil вместо &user
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User, updates map[string]interface{}) error {
	return r.db.Model(user).Updates(updates).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error // Удаляем запись из базы по ID
}

func (r *userRepository) FindByIDWithTx(tx *gorm.DB, id uint) (*models.User, error) {
	if tx == nil {
		return nil, errors.New("transaction is required")
	}

	var user models.User
	err := tx.First(&user, id).Error
	return &user, err
}

func (r *userRepository) UpdateWithTx(tx *gorm.DB, user *models.User) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	return tx.Save(user).Error
}
