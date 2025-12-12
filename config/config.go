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
	JWT_SECRET = os.Getenv("JWT_SECRET")
}
