package repository

import (
	"errors"

	"go-api/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByTitle(title string, userID uint) (*models.Category, error)
	GetByCategoryId(CategoryID uint, userID uint) (*models.Category, error)
	GetByCategoryIdWithTx(tx *gorm.DB, CategoryID uint, userID uint) (*models.Category, error)
	UpdateCategory(category *models.Category, updates map[string]interface{}) error
	UpdateWithTx(tx *gorm.DB, category *models.Category) error
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

func (c *categoryRepository) FindByTitle(title string, userID uint) (*models.Category, error) {
	var category models.Category
	err := c.db.Where("title = ?", title).Where("user_id = ?", userID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Возвращаем nil вместо &category
		}
		return nil, err // Возвращаем nil для других ошибок
	}
	return &category, nil
}

func (c *categoryRepository) GetByCategoryId(CategoryID uint, userID uint) (*models.Category, error) {
	var category models.Category
	err := c.db.Where("category_id =?", CategoryID).Where("user_id = ?", userID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &category, nil
}

func (c *categoryRepository) UpdateCategory(category *models.Category, updates map[string]interface{}) error {
	return c.db.Model(category).Updates(updates).Error
}

func (c *categoryRepository) GetByCategoryIdWithTx(tx *gorm.DB, CategoryID uint, userID uint) (*models.Category, error) {
	if tx == nil {
		return nil, errors.New("transaction is required")
	}

	var category models.Category
	err := tx.Where("category_id =?", CategoryID).Where("user_id = ?", userID).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) UpdateWithTx(tx *gorm.DB, category *models.Category) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	return tx.Save(category).Error // Используем GORM для обновления записи
}
