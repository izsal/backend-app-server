package main

import (
	"todo-app-backend/config"
	"todo-app-backend/db"
	"todo-app-backend/handler"
	"todo-app-backend/middleware"
	"todo-app-backend/repository"
	"todo-app-backend/service"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Init()
	db.Init()

	userRepo := repository.NewUserRepository(db.DB)
	todoRepo := repository.NewTodoRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)

	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	e := echo.New()

	// CORS middleware
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Logger middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	e.POST("/api/register", authHandler.Register)
	e.POST("/api/login", authHandler.Login)

	// Debug endpoint untuk test JWT
	e.GET("/api/debug", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "Public endpoint works"})
	})

	// Protected routes for todos
	todoGroup := e.Group("/api/todos")
	todoGroup.Use(middleware.JWTMiddleware)

	// Debug endpoint dengan JWT
	todoGroup.GET("/debug", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "JWT endpoint works", "user": c.Get("user")})
	})

	todoGroup.GET("", todoHandler.GetTodos)
	todoGroup.POST("", todoHandler.CreateTodo)
	todoGroup.PUT("/:id", todoHandler.UpdateTodo)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodo)

	// Protected routes for transactions
	transactionGroup := e.Group("/api/transactions")
	transactionGroup.Use(middleware.JWTMiddleware)

	transactionGroup.GET("", transactionHandler.GetTransactions)
	transactionGroup.GET("/summary", transactionHandler.GetTransactionSummary)
	transactionGroup.POST("", transactionHandler.CreateTransaction)
	transactionGroup.PUT("/:id", transactionHandler.UpdateTransaction)
	transactionGroup.DELETE("/:id", transactionHandler.DeleteTransaction)

	e.Logger.Fatal(e.Start(":8080"))
}
