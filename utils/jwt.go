package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetUserIDFromToken(c echo.Context) uint {
	// Ambil token dari context yang sudah di-set oleh JWT middleware
	tokenInterface := c.Get("user")
	if tokenInterface == nil {
		return 0
	}

	// Type assertion dengan safety check
	token, ok := tokenInterface.(*jwt.Token)
	if !ok {
		return 0
	}

	// Type assertion untuk claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}

	// Get user_id dari claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0
	}

	return uint(userIDFloat)
}
