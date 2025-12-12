package repository

import (
	"todo-app-backend/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByUserID(userID uint) ([]models.Transaction, error)
	FindByID(id uint) (*models.Transaction, error)
	Create(transaction *models.Transaction) error
	Update(transaction *models.Transaction) error
	Delete(transaction *models.Transaction) error
	GetTransactionSummary(userID uint) (map[string]float64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ?", userID).Order("date DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(transaction *models.Transaction) error {
	return r.db.Delete(transaction).Error
}

func (r *transactionRepository) GetTransactionSummary(userID uint) (map[string]float64, error) {
	var incomeTotal float64
	var expenseTotal float64

	// Get total income
	r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, models.Income).
		Select("COALESCE(SUM(amount), 0)").
		Row().
		Scan(&incomeTotal)

	// Get total expense
	r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, models.Expense).
		Select("COALESCE(SUM(amount), 0)").
		Row().
		Scan(&expenseTotal)

	return map[string]float64{
		"income":  incomeTotal,
		"expense": expenseTotal,
		"balance": incomeTotal - expenseTotal,
	}, nil
}
