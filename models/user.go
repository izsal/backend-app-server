package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null;size:50" json:"username" validate:"required,min=3,max=50"`
	Password     string `gorm:"not null;size:255" json:"-" validate:"required,min=6"`
	Todos        []Todo
	Transactions []Transaction
}
