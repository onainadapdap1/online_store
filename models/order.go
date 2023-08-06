package models

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserID            uint    `json:"user_id"`
	PaymentCategoryID uint    `json:"payment_category_id"`
	PaymentMethodID   uint    `json:"payment_method_id"`
	ReceiverName      string  `json:"receiver_name"`
	ProofOfPayment    string  `json:"proof_of_payment"`
	TotalPrice        float64 `json:"total_price"`
	Status            string  `json:"status"`
	User              User
	PaymentCategory   PaymentCategory
	PaymentMethod     PaymentMethod
	OrderItems        []OrderItem
}

func (o *Order) TableName() string {
	return "tb_orders"
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `json:"order_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
	Product     Product `json:"product"`
	Order       Order   `json:"order"`
}

func (oderItem *OrderItem) TableName() string {
	return "tb_order_items"
}
