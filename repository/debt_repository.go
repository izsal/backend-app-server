package repository

import (
	"todo-app-backend/models"

	"gorm.io/gorm"
)

type DebtRepository interface {
	FindByUserID(userID uint) ([]models.Debt, error)
	FindByID(id uint) (*models.Debt, error)
	Create(debt *models.Debt) error
	Update(debt *models.Debt) error
	Delete(debt *models.Debt) error
	FindByUserIDAndType(userID uint, debtType string) ([]models.Debt, error)
	GetDebtSummary(userID uint) (map[string]float64, error)
}

type debtRepository struct {
	db *gorm.DB
}

func NewDebtRepository(db *gorm.DB) DebtRepository {
	return &debtRepository{db}
}

func (r *debtRepository) FindByUserID(userID uint) ([]models.Debt, error) {
	var debts []models.Debt
	err := r.db.Where("user_id = ?", userID).Order("date DESC").Find(&debts).Error
	return debts, err
}

func (r *debtRepository) FindByUserIDAndType(userID uint, debtType string) ([]models.Debt, error) {
	var debts []models.Debt
	err := r.db.Where("user_id = ? AND type = ?", userID, debtType).Order("date DESC").Find(&debts).Error
	return debts, err
}

func (r *debtRepository) FindByID(id uint) (*models.Debt, error) {
	var debt models.Debt
	err := r.db.First(&debt, id).Error
	return &debt, err
}

func (r *debtRepository) Create(debt *models.Debt) error {
	return r.db.Create(debt).Error
}

func (r *debtRepository) Update(debt *models.Debt) error {
	return r.db.Save(debt).Error
}

func (r *debtRepository) Delete(debt *models.Debt) error {
	return r.db.Delete(debt).Error
}

func (r *debtRepository) GetDebtSummary(userID uint) (map[string]float64, error) {
	var totalOwed float64
	var totalOwing float64

	// Get total amount I am owed
	r.db.Model(&models.Debt{}).
		Where("user_id = ? AND type = ? AND status != ?", userID, models.DebtOwed, "paid").
		Select("COALESCE(SUM(remaining), 0)").
		Row().
		Scan(&totalOwed)

	// Get total amount I owe
	r.db.Model(&models.Debt{}).
		Where("user_id = ? AND type = ? AND status != ?", userID, models.DebtOwing, "paid").
		Select("COALESCE(SUM(remaining), 0)").
		Row().
		Scan(&totalOwing)

	return map[string]float64{
		"owed":    totalOwed,
		"owing":   totalOwing,
		"balance": totalOwed - totalOwing, // Positive means I'm owed more than I owe
	}, nil
}
