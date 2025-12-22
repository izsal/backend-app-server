package repository

import (
	"todo-app-backend/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByUserIDAndType(userID uint, categoryType string) ([]models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(category *models.Category) error
	FindByNameAndUserIDAndType(name string, userID uint, categoryType string) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindByUserIDAndType(userID uint, categoryType string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("user_id = ? AND type = ?", userID, categoryType).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(category *models.Category) error {
	return r.db.Delete(category).Error
}

func (r *categoryRepository) FindByNameAndUserIDAndType(name string, userID uint, categoryType string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ? AND user_id = ? AND type = ?", name, userID, categoryType).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
