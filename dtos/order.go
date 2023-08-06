package dtos

import "github.com/onainadapdap1/online_store/models"

type OrderFormatter struct {
	ID                uint                     `json:"id_order"`
	UserID            uint                     `json:"user_id"`
	User              OrderUserFormatter       `json:"user"`
	PaymentCategoryID uint                     `json:"payment_category_id"`
	PaymentCategory   OrderPaymentCatFormatter `json:"payment_category"`
	PaymentMethodID   uint                     `json:"payment_method_id"`
	PaymentMethod     OrderPaymentMetFormatter `json:"payment_method"`
	ReceiverName      string                   `json:"receiver_name"`
	ProofOfPayment    string                   `json:"proof_of_payment"`
	TotalPrice        float64                  `json:"total_price"`
	Status            string                   `json:"status"`
}
type OrderUserFormatter struct {
	FullName string `json:"full_name"`
}
type OrderPaymentCatFormatter struct {
	ID           uint   `json:"payment_category_id"`
	CategoryName string `json:"payment_category_name"`
}
type OrderPaymentMetFormatter struct {
	ID         uint   `json:"payment_method_id"`
	MethodName string `json:"payment_method_name"`
}

func FormatOrderDetail(order models.Order) OrderFormatter {
	orderFormatter := OrderFormatter{
		ID:                order.ID,
		UserID:            order.UserID,
		PaymentCategoryID: order.PaymentCategoryID,
		PaymentMethodID:   order.PaymentMethodID,
		ReceiverName:      order.ReceiverName,
		ProofOfPayment:    order.ProofOfPayment,
		TotalPrice:        order.TotalPrice,
		Status:            order.Status,
	}
	user := order.User
	orderUserFormatter := OrderUserFormatter{
		FullName: user.FullName,
	}
	orderFormatter.User = orderUserFormatter

	paymentCategory := order.PaymentCategory
	orderPaymentCategory := OrderPaymentCatFormatter{
		ID:           paymentCategory.ID,
		CategoryName: paymentCategory.CategoryName,
	}
	orderFormatter.PaymentCategory = orderPaymentCategory

	paymentMethod := order.PaymentMethod
	orderPaymentMethod := OrderPaymentMetFormatter{
		ID:         paymentMethod.ID,
		MethodName: paymentMethod.MethodName,
	}
	orderFormatter.PaymentMethod = orderPaymentMethod

	return orderFormatter

}
