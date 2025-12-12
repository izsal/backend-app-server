package handler

import (
	"net/http"
	"os"
	"time"
	"todo-app-backend/service"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service   service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		service:   s,
		validator: validator.New(),
	}
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var body AuthRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	// Validasi menggunakan validator
	if err := h.validator.Struct(body); err != nil {
		// Parse error dari validator untuk memberikan pesan yang lebih user-friendly
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				switch err.Tag() {
				case "required":
					validationErrors["username"] = "Username is required"
				case "min":
					validationErrors["username"] = "Username must be at least 3 characters"
				case "max":
					validationErrors["username"] = "Username must be at most 50 characters"
				}
			case "Password":
				switch err.Tag() {
				case "required":
					validationErrors["password"] = "Password is required"
				case "min":
					validationErrors["password"] = "Password must be at least 6 characters"
				}
			}
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validationErrors})
	}

	if err := h.service.Register(body.Username, body.Password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username already exists"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "Registered successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var body AuthRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	// Validasi menggunakan validator
	if err := h.validator.Struct(body); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				validationErrors["username"] = "Username is required"
			case "Password":
				validationErrors["password"] = "Password is required"
			}
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validationErrors})
	}

	user, err := h.service.Login(body.Username, body.Password)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	t, _ := token.SignedString([]byte(secret))
	return c.JSON(http.StatusOK, echo.Map{"token": t})
}
