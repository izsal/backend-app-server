package service

import (
	"errors"
	"time"
	"todo-app-backend/models"
	"todo-app-backend/repository"
)

type TransactionService interface {
	GetTransactions(userID uint) ([]models.Transaction, error)
	GetTransactionSummary(userID uint) (map[string]float64, error)
	CreateTransaction(userID uint, transactionReq CreateTransactionRequest) (*models.Transaction, error)
	UpdateTransaction(userID uint, id uint, transactionReq UpdateTransactionRequest) (*models.Transaction, error)
	DeleteTransaction(userID uint, id uint) error
}

type CreateTransactionRequest struct {
	Type        models.TransactionType `json:"type" validate:"required,oneof=income expense"`
	Category    string                 `json:"category" validate:"required,max=100"`
	Description string                 `json:"description" validate:"max=255"`
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Date        time.Time              `json:"date" validate:"required"`
}

type UpdateTransactionRequest struct {
	Type        models.TransactionType `json:"type" validate:"required,oneof=income expense"`
	Category    string                 `json:"category" validate:"required,max=100"`
	Description string                 `json:"description" validate:"max=255"`
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Date        time.Time              `json:"date" validate:"required"`
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

func (s *transactionService) GetTransactions(userID uint) ([]models.Transaction, error) {
	return s.repo.FindByUserID(userID)
}

func (s *transactionService) GetTransactionSummary(userID uint) (map[string]float64, error) {
	return s.repo.GetTransactionSummary(userID)
}

func (s *transactionService) CreateTransaction(userID uint, transactionReq CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		UserID:      userID,
		Type:        transactionReq.Type,
		Category:    transactionReq.Category,
		Description: transactionReq.Description,
		Amount:      transactionReq.Amount,
		Date:        transactionReq.Date,
	}

	err := s.repo.Create(transaction)
	return transaction, err
}

func (s *transactionService) UpdateTransaction(userID uint, id uint, transactionReq UpdateTransactionRequest) (*models.Transaction, error) {
	transaction, err := s.repo.FindByID(id)
	if err != nil || transaction.UserID != userID {
		return nil, errors.New("transaction not found")
	}

	transaction.Type = transactionReq.Type
	transaction.Category = transactionReq.Category
	transaction.Description = transactionReq.Description
	transaction.Amount = transactionReq.Amount
	transaction.Date = transactionReq.Date

	err = s.repo.Update(transaction)
	return transaction, err
}

func (s *transactionService) DeleteTransaction(userID uint, id uint) error {
	transaction, err := s.repo.FindByID(id)
	if err != nil || transaction.UserID != userID {
		return errors.New("transaction not found")
	}

	return s.repo.Delete(transaction)
}
