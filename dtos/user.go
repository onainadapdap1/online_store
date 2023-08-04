package dtos

import "github.com/onainadapdap1/online_store/models"

/* USER INPUT*/
// register user input field
type RegisterUserInput struct {
	FullName string `gorm:"not null" json:"full_name" form:"full_name" `
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" `
	Password string `gorm:"not null" json:"password" form:"password" `
}

type UserRegisterFormatter struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

/*USER*/
// format user response register
func FormatUserRegister(user *models.User) UserRegisterFormatter {
	formatter := UserRegisterFormatter{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
	}

	return formatter
}

type LoginUserInput struct {
	Email string `gorm:"not null;" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimun length of 6 characters"`
}

type UserLoginFormatter struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

func FormatUserLogin(user models.User, token string) UserLoginFormatter {
	userLoginFormatter := UserLoginFormatter {
		ID: user.ID,
		FullName: user.FullName,
		Email:  user.Email,
		Role: user.Role,
		Token: token,
	}

	return userLoginFormatter
}