package db

import (
	"log"
	"os"
	"todo-app-backend/config"
	"todo-app-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	// Check if we should use PostgreSQL or MySQL
	// Default to MySQL for backward compatibility
	usePostgres := os.Getenv("USE_POSTGRES") == "true"

	if usePostgres {
		// Use PostgreSQL
		DB, err = gorm.Open(postgres.Open(config.MYSQL_DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("PostgreSQL Connection failed: %v", err)
		}
	} else {
		// Use MySQL (default)
		DB, err = gorm.Open(mysql.Open(config.MYSQL_DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("MySQL Connection failed: %v", err)
		}
	}

	DB.AutoMigrate(&models.User{}, &models.Todo{}, &models.Transaction{}, &models.Category{}, &models.Debt{})
}
