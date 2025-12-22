package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name" validate:"required,max=100"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"size:20;not null;check:type IN ('income', 'expense')" json:"type" validate:"required,oneof=income expense"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}
