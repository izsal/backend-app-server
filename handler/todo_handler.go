package handler

import (
	"net/http"
	"strconv"
	"todo-app-backend/service"
	"todo-app-backend/utils"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	return &TodoHandler{s}
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	todos, err := h.service.GetTodos(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	var body struct {
		Title string `json:"title"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid"})
	}
	todo, err := h.service.CreateTodo(userID, body.Title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid"})
	}
	todo, err := h.service.UpdateTodo(userID, uint(id), body.Title, body.Completed)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Not found"})
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	userID := utils.GetUserIDFromToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteTodo(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Not found"})
	}
	return c.NoContent(http.StatusNoContent)
}