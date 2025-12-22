package service

import (
	"errors"
	"todo-app-backend/models"
	"todo-app-backend/repository"
)

type CategoryService interface {
	GetCategoriesByType(userID uint, categoryType string) ([]models.Category, error)
	CreateCategory(userID uint, categoryReq CreateCategoryRequest) (*models.Category, error)
	UpdateCategory(userID uint, id uint, categoryReq UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(userID uint, id uint) error
}

type CreateCategoryRequest struct {
	Name string                 `json:"name" validate:"required,max=100"`
	Type models.TransactionType `json:"type" validate:"required,oneof=income expense"`
}

type UpdateCategoryRequest struct {
	Name string                 `json:"name" validate:"required,max=100"`
	Type models.TransactionType `json:"type" validate:"required,oneof=income expense"`
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetCategoriesByType(userID uint, categoryType string) ([]models.Category, error) {
	return s.repo.FindByUserIDAndType(userID, categoryType)
}

func (s *categoryService) CreateCategory(userID uint, categoryReq CreateCategoryRequest) (*models.Category, error) {
	// Check if category with same name, user, and type already exists
	existingCategory, err := s.repo.FindByNameAndUserIDAndType(categoryReq.Name, userID, string(categoryReq.Type))
	if err == nil && existingCategory != nil {
		return nil, errors.New("category with this name already exists for this type")
	}

	category := &models.Category{
		Name:   categoryReq.Name,
		UserID: userID,
		Type:   string(categoryReq.Type),
	}

	err = s.repo.Create(category)
	return category, err
}

func (s *categoryService) UpdateCategory(userID uint, id uint, categoryReq UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil || category.UserID != userID {
		return nil, errors.New("category not found")
	}

	// Check if another category with same name, user, and type already exists
	existingCategory, err := s.repo.FindByNameAndUserIDAndType(categoryReq.Name, userID, string(categoryReq.Type))
	if err == nil && existingCategory != nil && existingCategory.ID != id {
		return nil, errors.New("category with this name already exists for this type")
	}

	category.Name = categoryReq.Name
	category.Type = string(categoryReq.Type)

	err = s.repo.Update(category)
	return category, err
}

func (s *categoryService) DeleteCategory(userID uint, id uint) error {
	category, err := s.repo.FindByID(id)
	if err != nil || category.UserID != userID {
		return errors.New("category not found")
	}

	return s.repo.Delete(category)
}
