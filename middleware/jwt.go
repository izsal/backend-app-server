package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"todo-app-backend/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Authorization header
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			fmt.Println("JWT middleware error: missing Authorization header")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or missing token",
			})
		}

		// Check Bearer prefix
		if !strings.HasPrefix(auth, "Bearer ") {
			fmt.Println("JWT middleware error: invalid Authorization header format")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or missing token",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == "" {
			fmt.Println("JWT middleware error: empty token")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or missing token",
			})
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWT_SECRET), nil
		})

		if err != nil {
			fmt.Printf("JWT middleware error: token parse failed: %v\n", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or missing token",
			})
		}

		if !token.Valid {
			fmt.Println("JWT middleware error: invalid token")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "Invalid or missing token",
			})
		}

		// Set token in context for handlers to use
		c.Set("user", token)

		return next(c)
	}
}
