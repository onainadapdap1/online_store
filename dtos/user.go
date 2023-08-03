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