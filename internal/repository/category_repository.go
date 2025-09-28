package repository

import (
	"errors"

	"go-api/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByTitle(title string, userID int) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (c *categoryRepository) Create(category *models.Category) error {
	return c.db.Create(category).Error
}

func (c *categoryRepository) FindByTitle(title string, userID int) (*models.Category, error) {
	var category models.Category
	err := c.db.Where("title = ?", title).Where("user_id = ?", userID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Возвращаем nil вместо &user
		}
		return nil, err // Возвращаем nil для других ошибок
	}
	return &category, nil
}
