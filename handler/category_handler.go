package handler

import (
	"net/http"
	"strconv"
	"todo-app-backend/service"
	"todo-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{s}
}

func (h *CategoryHandler) GetCategoriesByType(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	categoryType := c.QueryParam("type")

	if categoryType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "type parameter is required"})
	}

	categories, err := h.service.GetCategoriesByType(userID, categoryType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	var body service.CreateCategoryRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	category, err := h.service.CreateCategory(userID, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var body service.UpdateCategoryRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	category, err := h.service.UpdateCategory(userID, uint(id), body)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Category not found"})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.DeleteCategory(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Category not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
