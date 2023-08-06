package service

import (
	"errors"
	"fmt"

	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/helpers"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type UserServiceInterface interface {
	RegisterUser(input dtos.RegisterUserInput) (models.User, error)
	LoginUser(input dtos.LoginUserInput) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

type userService struct {
	repo repository.UserRepoInterface
}

func NewUserService(repo repository.UserRepoInterface) UserServiceInterface {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(input dtos.RegisterUserInput) (models.User, error) {
	hashedPassword, err := helpers.HassPass(input.Password)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		FullName: input.FullName,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	newUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *userService) LoginUser(input dtos.LoginUserInput) (models.User, error) {
	inputEmail := input.Email
	inputPassword := input.Password

	user, err := s.repo.FindByEmail(inputEmail)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on that email")
	}

	comparePass := helpers.ComparePassword([]byte(user.Password), []byte(inputPassword))
	if !comparePass {
		// fmt.Printf("Stored Password: %s\n", user.Password)
		// hashedInputPassword, _ := helpers.HassPass(inputPassword)
		// fmt.Printf("Hashed Input Password: %s\n", hashedInputPassword)
		return models.User{}, errors.New("invalid credentials") // Return a generic error message to prevent guessing valid passwords
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (models.User, error) {
	user, err := s.repo.GetUserByID(id)
	fmt.Println("error : ", err)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found with that id")
	}

	return user, nil
}