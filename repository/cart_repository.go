package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/online_store/models"
)

type CartRepositoryInterface interface {
	// FindCart(cartID uint) (models.Cart, error)
	CheckIsProductExist(productID uint) (models.Product, error)
	CheckIsUserHasACart(userID uint) (models.Cart, error)
	CreateCart(cart models.Cart) (models.Cart, error)
	CheckIsProductInCart(cartID uint, productID uint) (models.CartItem, error)
	CreateCartItem(cartItem models.CartItem) (models.CartItem, error)
	UpdateCartItem(cartItem models.CartItem) (models.CartItem, error)
	FindItem(cartID uint, productID uint) (models.CartItem, error)
	DeleteItem(item models.CartItem) error
	GetAllCartItems() ([]models.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepositoryInterface {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetAllCartItems() ([]models.CartItem, error) {
	tx := r.db.Begin()
	cartItems := []models.CartItem{}
	if err := tx.Debug().Preload("Product").Preload("Product.Category").Preload("Product.User").Preload("Cart").Find(&cartItems).Error; err != nil {
		return  cartItems, err
	}
	return cartItems, nil
}
func (r *cartRepository) CheckIsProductInCart(cartID uint, productID uint) (models.CartItem, error) {
	tx := r.db.Begin()
	var cartItem models.CartItem
	if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartID, productID).Preload("Product").Preload("Product.Category").Preload("Product.Category.User").Preload("Product.User").Preload("Cart").First(&cartItem).Error; err != nil {
		tx.Rollback()
		return cartItem, err
	}
	tx.Commit()

	return cartItem, nil
}

func (r *cartRepository) CreateCartItem(cartItem models.CartItem) (models.CartItem, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID).Create(&cartItem).Error; err != nil {
		tx.Rollback()
		return cartItem, err
	}
	if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID).Preload("Product").Preload("Product.Category").Preload("Product.User").Preload("Cart").First(&cartItem).Error; err != nil {
		tx.Rollback()
		return cartItem, err
	}

	tx.Commit()

	return cartItem, nil
}

func (r *cartRepository) UpdateCartItem(cartItem models.CartItem) (models.CartItem, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Save(&cartItem).Error; err != nil {
		tx.Rollback()
		return cartItem, err
	}
	if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID).Preload("Product").Preload("Product.Category").Preload("Product.User").Preload("Cart").First(&cartItem).Error; err != nil {
		tx.Rollback()
		return cartItem, err
	}
	tx.Commit()

	return cartItem, nil
}

func (r *cartRepository) CreateCart(cart models.Cart) (models.Cart, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&cart).Error; err != nil {
		tx.Rollback()
		return cart, nil
	}
	tx.Commit()

	return cart, nil
}
func (r *cartRepository) CheckIsUserHasACart(userID uint) (models.Cart, error) {
	tx := r.db.Begin()
	var cart models.Cart
	if err := tx.Debug().Where("user_id = ?", userID).First(&cart).Error; err != nil {
		tx.Rollback()
		return cart, err
	}
	tx.Commit()

	return cart, nil
}

func (r *cartRepository) CheckIsProductExist(productID uint) (models.Product, error) {
	tx := r.db.Begin()
	var product models.Product
	if err := tx.Debug().Where("id = ?", productID).First(&product).Error; err != nil {
		tx.Rollback()
		return product, err
	}
	tx.Commit()

	return product, nil
}

func (r *cartRepository) FindItem(cartID uint, productID uint) (models.CartItem, error) {
	tx := r.db.Begin()

	var item models.CartItem
	// if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartID, productID).Preload("Product").Preload("Product.Category").Preload("Product.Category.User").Preload("Product.User").Preload("Cart").First(&cartItem).Error; err != nil {

	if err := tx.Debug().Where("cart_id = ? AND product_id = ?", cartID, productID).Preload("Product").Preload("Product.Category").Preload("Product.User").Preload("Cart").First(&item).Error; err != nil {
		tx.Rollback()
		return item, err
	}
	tx.Commit()

	return item, nil
}

func (r *cartRepository) DeleteItem(item models.CartItem) error {
	tx := r.db.Begin()

	if err := tx.Debug().Unscoped().Delete(&item).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
