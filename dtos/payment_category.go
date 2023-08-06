package dtos

import "github.com/onainadapdap1/online_store/models"


type CreatePaymentCategoryInput struct {
	CategoryName string `gorm:"not null" json:"payment_category_name" form:"payment_category_name"`
	User         models.User `json:"-"`
}

type PaymentCategoryFormatter struct {
	ID           uint                         `json:"id"`
	UserID       uint                         `json:"user_id"`
	Slug         string                       `json:"slug"`
	CategoryName string                       `json:"payment_category_name"`
	User         PaymentCategoryUserFormatter `json:"user"`
}

type PaymentCategoryUserFormatter struct {
	ID       uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func FormateCreatePaymentCategory(paymentCategory models.PaymentCategory) PaymentCategoryFormatter {
	paymentCategoryFormatter := PaymentCategoryFormatter{
		ID:           paymentCategory.ID,
		UserID:       paymentCategory.UserID,
		Slug:         paymentCategory.Slug,
		CategoryName: paymentCategory.CategoryName,
	}
	user := paymentCategory.User
	paymentCategoryUserFormatter := PaymentCategoryUserFormatter{
		ID:       user.ID,
		FullName: user.FullName,
		Role:     user.Role,
	}
	paymentCategoryFormatter.User = paymentCategoryUserFormatter

	return paymentCategoryFormatter
}
