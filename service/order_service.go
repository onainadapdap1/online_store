package service

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type OrderServiceInterface interface {
	CheckIsUserHasACart(userID uint) (models.Cart, error)
	GetAllCartItems(cartID uint) ([]models.CartItem, error)
	CountPriceInCartItems(cartItems []models.CartItem) float64
	GetPaymenCID(paymentCID uint) (models.PaymentCategory, error)
	GetPaymentMID(paymentMID uint) (models.PaymentMethod, error)
	CreateOrder(order models.Order) (models.Order, error)
	CreateOrderItem(order models.Order, cartItems []models.CartItem) (int, error)
	UpdateCartAndDeleteItems(cartID, userID uint) error
}

type orderService struct {
	repo repository.OrderRepositoryInterface
}

func NewOrderService(repo repository.OrderRepositoryInterface) OrderServiceInterface {
	return &orderService{repo: repo}
}
func (s *orderService) UpdateCartAndDeleteItems(cartID, userID uint) error {
	if err := s.repo.UpdateCartUpdatedAt(userID); err != nil {
		return err
	}
	if err := s.repo.DeleteCartItemsByCartID(cartID); err != nil {
		return err
	}
	return nil
}
func (s *orderService) CreateOrderItem(order models.Order, cartItems []models.CartItem) (int, error) {
	for _, item := range cartItems {
		product := item.Product
		if product.Quantity < item.Quantity {
			return http.StatusBadRequest, errors.New("insufficient product quantity")
		}

		orderItem := models.OrderItem{
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			TotalPrice:  item.TotalPrice,
			Product:     item.Product,
		}
		err := s.repo.FindOrderIDAndProductIDInOrderItem(order.ID, product.ID)
		if err == nil {
			return http.StatusBadRequest, nil
		} else if err != gorm.ErrRecordNotFound {
			return http.StatusBadRequest, nil
		}

		err = s.repo.CreateOrderItem(orderItem)
		if err != nil {
			return http.StatusBadRequest, nil
		}

		product.Quantity -= item.Quantity
		err = s.repo.UpdateProduct(product)
		if err != nil {
			return http.StatusBadRequest, nil
		}
	}

	return http.StatusOK, nil
}
func (s *orderService) CreateOrder(order models.Order) (models.Order, error) {
	orderCreated, err := s.repo.CreateOrder(order)
	if err != nil {
		return orderCreated, err
	}
	return orderCreated, nil
}
func (s *orderService) GetPaymentMID(paymentMID uint) (models.PaymentMethod, error) {
	paymentMethodID, err := s.repo.GetPaymentMID(paymentMID)
	if err != nil {
		return paymentMethodID, err
	}
	return paymentMethodID, nil
}
func (s *orderService) GetPaymenCID(paymentCID uint) (models.PaymentCategory, error) {
	paymentCategoryID, err := s.repo.GetPaymenCID(paymentCID)
	if err != nil {
		return paymentCategoryID, nil
	}
	return paymentCategoryID, nil
}
func (s *orderService) CheckIsUserHasACart(userID uint) (models.Cart, error) {
	cart, err := s.repo.CheckIsUserHasACart(userID)
	if err != nil {
		return cart, err
	}
	return cart, nil
}

func (s *orderService) GetAllCartItems(cartID uint) ([]models.CartItem, error) {
	cartItems, err := s.repo.GetAllCartItems(cartID)
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

func (s *orderService) CountPriceInCartItems(cartItems []models.CartItem) float64 {
	var totalPrice float64
	for _, cartItem := range cartItems {
		totalPrice += cartItem.Product.Price * float64(cartItem.Quantity)
	}

	return totalPrice
}
