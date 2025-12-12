package handler

import (
	"net/http"
	"strconv"
	"todo-app-backend/service"
	"todo-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{s}
}

func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	transactions, err := h.service.GetTransactions(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionSummary(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	summary, err := h.service.GetTransactionSummary(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	var body service.CreateTransactionRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	transaction, err := h.service.CreateTransaction(userID, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var body service.UpdateTransactionRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	transaction, err := h.service.UpdateTransaction(userID, uint(id), body)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Transaction not found"})
	}
	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteTransaction(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Transaction not found"})
	}
	return c.NoContent(http.StatusNoContent)
}
