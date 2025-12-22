package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      uint            `gorm:"not null" json:"user_id"`
	Type        TransactionType `gorm:"type:enum('income','expense');not null" json:"type" validate:"required,oneof=income expense"`
	CategoryID  *uint           `gorm:"index" json:"category_id,omitempty"`
	Category    string          `gorm:"size:100;not null" json:"category" validate:"required,max=100"`
	Description string          `gorm:"size:255" json:"description" validate:"max=255"`
	Amount      float64         `gorm:"not null" json:"amount" validate:"required,gt=0"`
	Date        time.Time       `gorm:"not null" json:"date" validate:"required"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User        User     `gorm:"foreignKey:UserID" json:"-"`
	CategoryObj Category `gorm:"foreignKey:CategoryID" json:"category_obj,omitempty"`
}
