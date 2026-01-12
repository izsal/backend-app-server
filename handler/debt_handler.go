package handler

import (
	"net/http"
	"strconv"
	"todo-app-backend/service"
	"todo-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type DebtHandler struct {
	service service.DebtService
}

func NewDebtHandler(s service.DebtService) *DebtHandler {
	return &DebtHandler{s}
}

func (h *DebtHandler) GetDebts(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	debts, err := h.service.GetDebts(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, debts)
}

func (h *DebtHandler) GetDebtsByType(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	debtType := c.QueryParam("type")

	if debtType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "type parameter is required"})
	}

	debts, err := h.service.GetDebtsByType(userID, debtType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, debts)
}

func (h *DebtHandler) GetDebtSummary(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	summary, err := h.service.GetDebtSummary(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}

func (h *DebtHandler) CreateDebt(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	var body service.CreateDebtRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	debt, err := h.service.CreateDebt(userID, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, debt)
}

func (h *DebtHandler) UpdateDebt(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var body service.UpdateDebtRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	debt, err := h.service.UpdateDebt(userID, uint(id), body)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Debt not found"})
	}
	return c.JSON(http.StatusOK, debt)
}

func (h *DebtHandler) DeleteDebt(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.DeleteDebt(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Debt not found"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *DebtHandler) MakePayment(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var body service.MakePaymentRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	debt, err := h.service.MakePayment(userID, uint(id), body)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Debt not found"})
	}
	return c.JSON(http.StatusOK, debt)
}
