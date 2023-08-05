package dtos

import "github.com/onainadapdap1/online_store/models"

type CreateCartItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
	User      models.User
	Product   models.Product
	Cart      models.Cart
}

// formatter create cart
type CartItemFormatter struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	CartID     uint    `json:"cart_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	Product    CartProductFormatter
	Cart       CartItemCartFormatter
}
type CartItemCartFormatter struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}
type CartProductFormatter struct {
	ID          uint    `json:"product_id"`
	Name        string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    CartProductCategoryFormatter
}
type CartProductCategoryFormatter struct {
	ID          uint   `json:"category_id"`
	Name        string `json:"product_name"`
	Description string `json:"description"`
	User        CartProductCategoryUserFormatter
}
type CartProductCategoryUserFormatter struct {
	ID       uint   `json:"user_id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func FormatListCartItems(cartItems []models.CartItem) []CartItemFormatter {
	cartItemsFormatter := []CartItemFormatter{}
	for _, cartItem := range cartItems {
		cartItemFormatter := FormatCreateCart(cartItem)
		cartItemsFormatter = append(cartItemsFormatter, cartItemFormatter)
	}
	return cartItemsFormatter
}

func FormatCreateCart(cartItem models.CartItem) CartItemFormatter {
	cartItemFormatter := CartItemFormatter{
		ID:         cartItem.ID,
		UserID:     cartItem.Cart.UserID,
		CartID:     cartItem.CartID,
		ProductID:  cartItem.ProductID,
		Quantity:   cartItem.Quantity,
		Price:      cartItem.Price,
		TotalPrice: cartItem.TotalPrice,
	}
	user := cartItem.Product.User
	categoryUser := CartProductCategoryUserFormatter{
		ID:       user.ID,
		FullName: user.FullName,
		Role:     user.Role,
	}

	categoryProduct := cartItem.Product.Category
	category := CartProductCategoryFormatter{
		ID:          categoryProduct.ID,
		Name:        categoryProduct.Name,
		Description: categoryProduct.Description,
		User:        categoryUser,
	}

	product := cartItem.Product
	prod := CartProductFormatter{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    category,
	}

	cartItemFormatter.Product = prod

	cart := cartItem.Cart
	cartData := CartItemCartFormatter{
		ID:     cart.ID,
		UserID: cart.UserID,
	}

	cartItemFormatter.Cart = cartData

	return cartItemFormatter
}
