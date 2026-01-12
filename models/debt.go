package models

import (
	"time"

	"gorm.io/gorm"
)

type DebtType string

const (
	DebtOwed  DebtType = "owed"  // Debt someone owes me
	DebtOwing DebtType = "owing" // Debt I owe someone
)

type Debt struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Type        DebtType       `gorm:"type:enum('owed','owing');not null" json:"type" validate:"required,oneof=owed owing"`
	Name        string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"` // Person/Entity name
	Description string         `gorm:"size:255" json:"description" validate:"max=255"`
	Amount      float64        `gorm:"not null" json:"amount" validate:"required,gt=0"`
	Remaining   float64        `gorm:"not null;default:0" json:"remaining" validate:"gte=0"` // Amount still owed
	Date        time.Time      `gorm:"not null" json:"date" validate:"required"`
	DueDate     *time.Time     `gorm:"" json:"due_date,omitempty"`                                                                    // Optional due date
	Status      string         `gorm:"size:20;not null;default:'active'" json:"status" validate:"required,oneof=active paid partial"` // active, paid, partial
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}
