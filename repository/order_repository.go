package repository

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type OrderRepositoryInterface interface {
	CheckIsUserHasACart(userID uint) (models.Cart, error)
	GetAllCartItems(cartID uint) ([]models.CartItem, error)
	GetPaymenCID(paymentCID uint) (models.PaymentCategory, error)
	GetPaymentMID(paymentMID uint) (models.PaymentMethod, error)
	CreateOrder(order models.Order) (models.Order, error)
	CreateOrderItem(orderItem models.OrderItem) error
	FindOrderIDAndProductIDInOrderItem(orderID, productID uint) error
	UpdateProduct(product models.Product) error
	UpdateCartUpdatedAt(userID uint) error
	DeleteCartItemsByCartID(cartID uint) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &orderRepository{db: db}
}

func (r *orderRepository) DeleteCartItemsByCartID(cartID uint) error {
	tx := r.db.Begin()
	if err := tx.Debug().Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (r *orderRepository) UpdateCartUpdatedAt(userID uint) error {
	tx := r.db.Begin()
	if err := tx.Debug().Model(&models.Cart{UserID: userID}).Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (r *orderRepository) UpdateProduct(product models.Product) error {
	tx := r.db.Begin()
	if err := tx.Debug().Save(&product).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (r *orderRepository) CreateOrderItem(orderItem models.OrderItem) error {
	tx := r.db.Begin()
	if err := tx.Debug().Create(&orderItem).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *orderRepository) FindOrderIDAndProductIDInOrderItem(orderID, productID uint) error {
	tx := r.db.Begin()
	if err := tx.Debug().Where("order_id = ? AND product_id = ?", orderID, productID).First(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *orderRepository) CreateOrder(order models.Order) (models.Order, error) {
	tx := r.db.Begin()
	if err := tx.Debug().Create(&order).Error; err != nil {
		tx.Rollback()
		return order, err
	}
	// Preload the User association and append it to the Order
	var user models.User
	if err := r.db.Debug().First(&user, order.UserID).Error; err != nil {
		tx.Rollback()
		return order, err
	}

	// Assign the fetched User to the order's User field
	order.User = user

	tx.Commit()

	return order, nil
}

func (r *orderRepository) CheckIsUserHasACart(userID uint) (models.Cart, error) {
	tx := r.db.Begin()
	var cart models.Cart
	if err := tx.Debug().Where("user_id = ?", userID).First(&cart).Error; err != nil {
		tx.Rollback()
		return cart, err
	}
	tx.Commit()

	return cart, nil
}

func (r *orderRepository) GetAllCartItems(cartID uint) ([]models.CartItem, error) {
	tx := r.db.Begin()
	var cartItems []models.CartItem
	if err := tx.Debug().Where("cart_id = ?", cartID).Preload("Product").Find(&cartItems).Error; err != nil {
		tx.Rollback()
		return cartItems, err
	}
	tx.Commit()

	return cartItems, nil
}
func (r *orderRepository) GetPaymenCID(paymentCID uint) (models.PaymentCategory, error) {
	tx := r.db.Begin()
	var paymentCategoryID models.PaymentCategory
	if err := tx.Debug().Where("id = ?", paymentCID).First(&paymentCategoryID).Error; err != nil {
		tx.Rollback()
		return paymentCategoryID, err
	}
	tx.Commit()
	return paymentCategoryID, nil
}

func (r *orderRepository) GetPaymentMID(paymentMID uint) (models.PaymentMethod, error) {
	tx := r.db.Begin()
	var paymentMethodID models.PaymentMethod
	if err := tx.Debug().Where("id = ?", paymentMID).First(&paymentMethodID).Error; err != nil {
		tx.Rollback()
		return paymentMethodID, err
	}
	tx.Commit()
	return paymentMethodID, nil
}
