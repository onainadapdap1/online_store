package models

import (
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/helpers"
)

// generate new user table
type User struct {
	gorm.Model
	FullName       string `gorm:"not null" json:"full_name" form:"full_name" valid:"required~Your Full Name is required"`
	Email          string `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your Email is required,email~Invalid Email format"`
	Password       string `gorm:"not null" json:"password" form:"password" valid:"required~Your Password is required,minstringlength(6)~Your Password must be at least 6 characters"`
	Role           string `json:"role"`
}

// naming convention
func (u *User) TableName() string {
	return "tb_users"
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password, err = helpers.HassPass(u.Password)
	if err != nil {
		log.Println("error while hashing password")
		return
	}
	err = nil
	return
}
