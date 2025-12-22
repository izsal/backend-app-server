package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MYSQL_DSN  string
	JWT_SECRET string
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	MYSQL_DSN = os.Getenv("MYSQL_DSN")
	if MYSQL_DSN == "" {
		MYSQL_DSN = "finance_user:finance_password@tcp(localhost:3306)/finance_app?charset=utf8mb4&parseTime=True&loc=Local"
	}

	JWT_SECRET = os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		JWT_SECRET = "your_jwt_secret_key_here"
	}
}
