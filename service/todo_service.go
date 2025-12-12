package service

import (
	"todo-app-backend/models"
	"todo-app-backend/repository"
)

type TodoService interface {
	GetTodos(userID uint) ([]models.Todo, error)
	CreateTodo(userID uint, title string) (*models.Todo, error)
	UpdateTodo(userID uint, id uint, title string, completed bool) (*models.Todo, error)
	DeleteTodo(userID uint, id uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo}
}

func (s *todoService) GetTodos(userID uint) ([]models.Todo, error) {
	return s.repo.FindByUserID(userID)
}

func (s *todoService) CreateTodo(userID uint, title string) (*models.Todo, error) {
	todo := &models.Todo{Title: title, UserID: userID}
	err := s.repo.Create(todo)
	return todo, err
}

func (s *todoService) UpdateTodo(userID uint, id uint, title string, completed bool) (*models.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if err != nil || todo.UserID != userID {
		return nil, err
	}
	todo.Title = title
	todo.Completed = completed
	err = s.repo.Update(todo)
	return todo, err
}

func (s *todoService) DeleteTodo(userID uint, id uint) error {
	todo, err := s.repo.FindByID(id)
	if err != nil || todo.UserID != userID {
		return err
	}
	return s.repo.Delete(todo)
}
