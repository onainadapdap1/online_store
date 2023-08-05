package models

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	UserID uint `json:"user_id"`
}

// CartItem model
type CartItem struct {
	gorm.Model
	CartID     uint    `json:"cart_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	Product    Product `gorm:"foreignkey:ProductID" json:"product"`
	Cart       Cart    `gorm:"foreignkey:CartID" json:"cart"`
}

func (c *Cart) TableName() string {
	return "tb_carts"
}

func (cartItem *CartItem) TableName() string {
	return "tb_cart_items"
}
