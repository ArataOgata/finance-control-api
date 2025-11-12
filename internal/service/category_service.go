package service

import (
	"errors"
	"fmt"
	dto "go-api/internal/dto/category_dto"
	"go-api/internal/models"
	"go-api/internal/repository"
	"log"

	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(title string, description string, userId uint) (*dto.CategoryResponse, error)
	GetCategory(CategoryID uint, userID uint) (*dto.CategoryResponse, error)
	UpdateCategory(req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
}

type categoryService struct {
	repo     repository.CategoryRepository
	userRepo repository.UserRepository
}

func NewCategoryService(repo repository.CategoryRepository, userRepo repository.UserRepository) CategoryService {
	return &categoryService{repo: repo, userRepo: userRepo}
}

func (c *categoryService) CreateCategory(title string, description string, userID uint) (*dto.CategoryResponse, error) {
	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User not found:", userID)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Println("Failed to check user:", err)
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	exists, err := c.repo.FindByTitle(title, user.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Пользователь не найден, можно продолжать регистрацию
		} else {
			log.Println("Failed to check title:", err)                        // Логируем ошибку
			return nil, fmt.Errorf("failed to check category title: %w", err) // Возвращаем обёрнутую ошибку
		}
	}

	if exists != nil {
		log.Println("Title already taken:", title)
		return nil, errors.New("title already taken")
	}

	category := &models.Category{
		Title:       title,
		Description: description,
		UserID:      user.UserID,
	}

	err = c.repo.Create(category)
	return &dto.CategoryResponse{
		UserID:      category.UserID,
		CategoryID:  category.CategoryID,
		Title:       category.Title,
		Description: &category.Title,
		Total:       category.Total,
	}, err
}

func (c *categoryService) GetCategory(CategoryID uint, userID uint) (*dto.CategoryResponse, error) {
	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User not found:", userID)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		log.Println("Failed to check user:", err)
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	category, err := c.repo.GetByCategoryId(CategoryID, user.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Category not found:", CategoryID)
			return nil, err
		}
		log.Println("Failed to get category:", err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &dto.CategoryResponse{
		UserID:      category.UserID,
		CategoryID:  category.CategoryID,
		Title:       category.Title,
		Description: &category.Title,
		Total:       category.Total,
	}, nil
}

func (c *categoryService) UpdateCategory(req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

	user, err := c.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	category, err := c.repo.GetByCategoryId(req.CategoryID, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if req.IsEmpty() {
		return nil, errors.New("no fields to update")
	}

	updates := req.ToMap()

	if err := c.repo.UpdateCategory(category, updates); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	if title, ok := updates["title"].(string); ok {
		category.Title = title
	}

	if description, ok := updates["description"].(string); ok {
		category.Description = description
	}

	if total, ok := updates["total"].(int); ok {
		category.Total = total
	}

	return &dto.CategoryResponse{
		UserID:      category.UserID,
		CategoryID:  category.CategoryID,
		Title:       category.Title,
		Description: &category.Title,
		Total:       category.Total,
	}, nil
}
