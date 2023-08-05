package dtos

import "github.com/onainadapdap1/online_store/models"

type CreatePaymentMethodInput struct {
	CategoryPaymentID uint   `gorm:"not null" form:"category_payment_id" json:"category_payment_id"`
	MethodName        string `gorm:"not null" form:"method_name" json:"method_name"`
	Number            string `gorm:"not null" form:"number" json:"number"`
	OwnerName         string `gorm:"not null" form:"owner_name" json:"owner_name"`
	User              models.User
	PaymentCategory   models.PaymentCategory
}

// func main() {
// 	pm := models.PaymentMethod{}
// 	pc := models.PaymentCategory{}
// }

// type PaymentMethod struct {
// 	gorm.Model
// 	CategoryPaymentID uint `gorm:"not null"`
// 	MethodName string `gorm:"not null"`
// 	Number string `json:"number"`
// 	OwnerName string `json:"owner_name"`
// 	CategoryName string
// 	PaymentCategory PaymentCategory `gorm:"foreignkey:CategoryPaymentID"`
// }

type PaymentMethodFormatter struct {
	ID                uint   `json:"id"`
	UserID            uint   `json:"user_id"`
	CategoryPaymentID uint   `json:"payment_category_id"`
	MethodName        string `json:"method_name"`
	Number            string `json:"number"`
	OwnerName         string `json:"owner_name"`
	CategoryName      string `json:"payment_category_name"`
	PaymentCategory   PaymentMethodCategoryFormatter
	User              PaymentMethodUserFormatter
}

type PaymentMethodCategoryFormatter struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"payment_category_name"`
}

type PaymentMethodUserFormatter struct {
	ID       uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func FormateCreatePaymentMethod(paymentMethod models.PaymentMethod) PaymentMethodFormatter {
	paymentMethodFormatter := PaymentMethodFormatter{
		ID:                paymentMethod.ID,
		UserID:            paymentMethod.UserID,
		CategoryPaymentID: paymentMethod.CategoryPaymentID,
		MethodName:        paymentMethod.MethodName,
		Number:            paymentMethod.Number,
		OwnerName:         paymentMethod.OwnerName,
		CategoryName:      paymentMethod.PaymentCategory.CategoryName,
	}
	paymentCategory := paymentMethod.PaymentCategory
	paymentCategoryFormatter := PaymentMethodCategoryFormatter{
		ID:           paymentCategory.ID,
		CategoryName: paymentCategory.CategoryName,
	}
	user := paymentMethod.User
	paymentUserFormatter := PaymentMethodUserFormatter{
		ID:       user.ID,
		FullName: user.FullName,
		Role:     user.Role,
	}
	paymentMethodFormatter.PaymentCategory = paymentCategoryFormatter
	paymentMethodFormatter.User = paymentUserFormatter

	return paymentMethodFormatter
}
