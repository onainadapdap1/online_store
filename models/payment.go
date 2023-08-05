package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type PaymentCategory struct {
	gorm.Model
	UserID uint 
	CategoryName string `gorm:"unique;not null"`
	Slug string `gorm:"unique;not null"`
	User User `gorm:"foreignkey:UserID"`
}

func (pc *PaymentCategory) TableName() string {
	return "tb_payment_categories"
}

func (pc *PaymentCategory) BeforeSave() (err error) {
	pc.Slug = slug.Make(pc.CategoryName)
	return
}

type PaymentMethod struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	CategoryPaymentID uint `gorm:"not null"`
	MethodName string `gorm:"not null"`
	Number string `json:"number"`
	OwnerName string `json:"owner_name"`
	CategoryName string
	User User
	PaymentCategory PaymentCategory `gorm:"foreignkey:CategoryPaymentID"`
}

func (pm *PaymentMethod) TableName() string {
	return "tb_payment_methods"
}
