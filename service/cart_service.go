package service

import (
	"log"

	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/repository"
)

type CartServiceInterface interface {
	CheckIsProductExist(productID uint) (models.Product, error)
	CheckIsUserHasACart(userID uint) (models.Cart, error)
	CreateCartItem(input dtos.CreateCartItemInput) (models.CartItem, error)
	FindItem(cartID uint, productID uint) (models.CartItem, error)
	UpdateCartItem(cartItem models.CartItem) (models.CartItem, error)
	DeleteItem(item models.CartItem) error
	GetAllUserCartItems(userCartID uint) ([]models.CartItem, error)
	GetCartByUserID(userID uint) (models.Cart, error)
}

type cartService struct {
	repo repository.CartRepositoryInterface
}

func NewCartService(repo repository.CartRepositoryInterface) CartServiceInterface {
	return &cartService{repo: repo}
}
func (s *cartService) GetCartByUserID(userID uint) (models.Cart, error) {
	userCart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return userCart, err
	}
	return userCart, nil
}
func (s *cartService) GetAllUserCartItems(userCartID uint) ([]models.CartItem, error) {
	cartItems, err := s.repo.GetAllUserCartItems(userCartID)
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}
func (s *cartService) UpdateCartItem(cartItem models.CartItem) (models.CartItem, error) {
	cartItem, err := s.repo.UpdateCartItem(cartItem)
	if err != nil {
		return cartItem, err
	}

	return cartItem, nil
}

func (s *cartService) CreateCartItem(input dtos.CreateCartItemInput) (models.CartItem, error) {
	// Check if the product is already in the cart
	if cartItem, err := s.repo.CheckIsProductInCart(input.Cart.ID, input.ProductID); err != nil {
		// If the product is not in the cart, create a new cart item
		cartItem := models.CartItem{
			CartID:     input.Cart.ID,
			ProductID:  input.ProductID,
			Quantity:   input.Quantity,
			Price:      input.Product.Price,
			TotalPrice: float64(input.Quantity) * input.Product.Price,
		}
		if createCartItem, err := s.repo.CreateCartItem(cartItem); err != nil {
			return createCartItem, err
		} else {
			return createCartItem, nil
		}
	} else {
		// If the product is already in the cart, update the quantity and total price
		cartItem.Quantity += input.Quantity
		cartItem.TotalPrice += float64(input.Quantity) * input.Product.Price
		if cartItem, err := s.repo.UpdateCartItem(cartItem); err != nil {
			return cartItem, err
		} else {
			return cartItem, nil
		}
	}
}

func (s *cartService) CheckIsUserHasACart(userID uint) (models.Cart, error) {
	cart, err := s.repo.CheckIsUserHasACart(userID)
	if err != nil {
		cart = models.Cart{
			UserID: userID,
		}
		if createdCart, err := s.repo.CreateCart(cart); err != nil {
			log.Println("cart is : ", cart)
			return cart, err
		} else {
			cart = createdCart
		}
	}

	return cart, nil
}

func (s *cartService) CheckIsProductExist(productID uint) (models.Product, error) {
	product, err := s.repo.CheckIsProductExist(productID)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *cartService) FindItem(cartID uint, productID uint) (models.CartItem, error) {
	item, err := s.repo.FindItem(cartID, productID)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (s *cartService) DeleteItem(item models.CartItem) error {
	err := s.repo.DeleteItem(item)
	if err != nil {
		return err
	}

	return nil
}
