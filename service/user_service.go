package service

import (
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type UserServiceInterface interface {
	RegisterUser(input dtos.RegisterUserInput) (models.User, error)
}

type userService struct {
	repo repository.UserRepoInterface
}

func NewUserService(repo repository.UserRepoInterface) UserServiceInterface {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(input dtos.RegisterUserInput) (models.User, error) {
	user := models.User {
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password,
		Role: "user",
	}

	newUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}