package service

import (
	"errors"
	"time"
	"todo-app-backend/models"
	"todo-app-backend/repository"
)

type DebtService interface {
	GetDebts(userID uint) ([]models.Debt, error)
	GetDebtsByType(userID uint, debtType string) ([]models.Debt, error)
	GetDebtSummary(userID uint) (map[string]float64, error)
	CreateDebt(userID uint, debtReq CreateDebtRequest) (*models.Debt, error)
	UpdateDebt(userID uint, id uint, debtReq UpdateDebtRequest) (*models.Debt, error)
	DeleteDebt(userID uint, id uint) error
	MakePayment(userID uint, id uint, paymentReq MakePaymentRequest) (*models.Debt, error)
}

type CreateDebtRequest struct {
	Type        models.DebtType `json:"type" validate:"required,oneof=owed owing"`
	Name        string          `json:"name" validate:"required,max=100"`
	Description string          `json:"description" validate:"max=255"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Date        time.Time       `json:"date" validate:"required"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
}

type UpdateDebtRequest struct {
	Type        models.DebtType `json:"type" validate:"required,oneof=owed owing"`
	Name        string          `json:"name" validate:"required,max=100"`
	Description string          `json:"description" validate:"max=255"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	Date        time.Time       `json:"date" validate:"required"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
}

type MakePaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type debtService struct {
	repo repository.DebtRepository
}

func NewDebtService(repo repository.DebtRepository) DebtService {
	return &debtService{repo}
}

func (s *debtService) GetDebts(userID uint) ([]models.Debt, error) {
	return s.repo.FindByUserID(userID)
}

func (s *debtService) GetDebtsByType(userID uint, debtType string) ([]models.Debt, error) {
	return s.repo.FindByUserIDAndType(userID, debtType)
}

func (s *debtService) GetDebtSummary(userID uint) (map[string]float64, error) {
	return s.repo.GetDebtSummary(userID)
}

func (s *debtService) CreateDebt(userID uint, debtReq CreateDebtRequest) (*models.Debt, error) {
	debt := &models.Debt{
		UserID:      userID,
		Type:        debtReq.Type,
		Name:        debtReq.Name,
		Description: debtReq.Description,
		Amount:      debtReq.Amount,
		Remaining:   debtReq.Amount, // Initially remaining equals the full amount
		Date:        debtReq.Date,
		DueDate:     debtReq.DueDate,
		Status:      "active", // Default to active
	}

	err := s.repo.Create(debt)
	return debt, err
}

func (s *debtService) UpdateDebt(userID uint, id uint, debtReq UpdateDebtRequest) (*models.Debt, error) {
	debt, err := s.repo.FindByID(id)
	if err != nil || debt.UserID != userID {
		return nil, errors.New("debt not found")
	}

	debt.Type = debtReq.Type
	debt.Name = debtReq.Name
	debt.Description = debtReq.Description
	debt.Amount = debtReq.Amount
	debt.Date = debtReq.Date
	debt.DueDate = debtReq.DueDate

	// If amount changed, update remaining accordingly
	if debtReq.Amount != debt.Amount {
		// Calculate the difference and adjust remaining
		paidAmount := debt.Amount - debt.Remaining
		debt.Remaining = debtReq.Amount - paidAmount
		if debt.Remaining <= 0 {
			debt.Status = "paid"
			debt.Remaining = 0
		} else if debt.Remaining < debtReq.Amount {
			debt.Status = "partial"
		} else {
			debt.Status = "active"
		}
	}

	err = s.repo.Update(debt)
	return debt, err
}

func (s *debtService) DeleteDebt(userID uint, id uint) error {
	debt, err := s.repo.FindByID(id)
	if err != nil || debt.UserID != userID {
		return errors.New("debt not found")
	}

	return s.repo.Delete(debt)
}

func (s *debtService) MakePayment(userID uint, id uint, paymentReq MakePaymentRequest) (*models.Debt, error) {
	debt, err := s.repo.FindByID(id)
	if err != nil || debt.UserID != userID {
		return nil, errors.New("debt not found")
	}

	// Calculate new remaining amount
	newRemaining := debt.Remaining - paymentReq.Amount
	if newRemaining < 0 {
		newRemaining = 0
	}

	debt.Remaining = newRemaining

	// Update status based on remaining amount
	if newRemaining <= 0 {
		debt.Status = "paid"
		debt.Remaining = 0
	} else if newRemaining < debt.Amount {
		debt.Status = "partial"
	} else {
		debt.Status = "active"
	}

	err = s.repo.Update(debt)
	return debt, err
}
