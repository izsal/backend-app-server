package db

import (
	"log"
	"todo-app-backend/config"
	"todo-app-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.MYSQL_DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB Connection failed: %v", err)
	}
	DB.AutoMigrate(&models.User{}, &models.Todo{}, &models.Transaction{}, &models.Category{})
}
