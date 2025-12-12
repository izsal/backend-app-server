package service

import (
	"todo-app-backend/models"
	"todo-app-backend/repository"
	"todo-app-backend/utils"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(username, password string) error {
	hash, _ := utils.HashPassword(password)
	user := &models.User{Username: username, Password: hash}
	return s.userRepo.Create(user)
}

func (s *authService) Login(username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, nil
	}
	return user, nil
}
