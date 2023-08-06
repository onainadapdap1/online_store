package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type CartHandlerInterface interface {
	DeleteCartItem(c *gin.Context)
	AddProductToCart(c *gin.Context)
	UpdateCartItemQuantity(c *gin.Context)
	GetAllUserCartItems(c *gin.Context)
}

type cartHandler struct {
	service service.CartServiceInterface
}

func NewCartHandler(service service.CartServiceInterface) CartHandlerInterface {
	return &cartHandler{service: service}
}

func (h *cartHandler) GetAllUserCartItems(c *gin.Context) {
	currentUser := c.MustGet("currentuser").(models.User)
	userID := currentUser.ID
	userCart, err := h.service.GetCartByUserID(userID)
	if err != nil {
		response := utils.APIResponse("failed to get user cart", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	cartItems, err := h.service.GetAllUserCartItems(userCart.ID)
	if err != nil {
		response := utils.APIResponse("failed to get all user's cart items", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("list of cart items", http.StatusOK, "success", dtos.FormatListCartItems(cartItems))
	c.JSON(http.StatusOK, response)
}

func (h *cartHandler) AddProductToCart(c *gin.Context) {
	// Get the current user ID from the request context
	currentUser := c.MustGet("currentuser").(models.User)
	userID := currentUser.ID
	var input dtos.CreateCartItemInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product exists
	var product models.Product
	product, err := h.service.CheckIsProductExist(input.ProductID)
	if err != nil {
		response := utils.APIResponse("Failed to find product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// cek apakah input lebih besar dari jumlah barang tersedia
	if input.Quantity > product.Quantity {
		response := utils.APIResponse("input quantity cannot be greater than stock", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Check if the user has a cart
	cart, err := h.service.CheckIsUserHasACart(userID)
	log.Println("cart in handler is : ", cart)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := utils.APIResponse("user cart checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Check if the product is already in the cart
	input.User = currentUser
	input.Product = product
	input.Cart = cart

	cartItem, err := h.service.CreateCartItem(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := utils.APIResponse("failed to create cart item", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := utils.APIResponse("success create cart item", http.StatusOK, "success", dtos.FormatCreateCart(cartItem))
	c.JSON(http.StatusOK, response)

}

// UpdateCartItemQuantity is a handler function to update the quantity of a product in a cart
func (h *cartHandler) UpdateCartItemQuantity(c *gin.Context) {
	// Parse cart and item IDs from the URL parameters
	cartID, _ := strconv.Atoi(c.Param("cart_id"))
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	quantity := c.PostForm("quantity")
	quantityUpdate, _ := strconv.Atoi(quantity)

	item, err := h.service.FindItem(uint(cartID), uint(itemID))
	if err != nil {
		response := utils.APIResponse("Failed find product item from cart item", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Calculate the new quantity for the cart item
	newQuantity := item.Quantity
	if c.PostForm("action") == "add" {
		newQuantity += quantityUpdate
	} else if c.PostForm("action") == "remove" {
		newQuantity -= quantityUpdate
	}

	// Check if the new quantity is valid
	if newQuantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Update the quantity of the cart item
	item.Quantity = newQuantity
	item.TotalPrice = float64(item.Quantity) * item.Price
	cartItem, err := h.service.UpdateCartItem(item)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := utils.APIResponse("failed to update cart item", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Return the updated cart item as JSON
	response := utils.APIResponse("success update cart item", http.StatusOK, "success", dtos.FormatCreateCart(cartItem))
	c.JSON(http.StatusOK, response)
}

func (h *cartHandler) DeleteCartItem(c *gin.Context) {
	// parse cart and item ids from the url parameters
	cartID, _ := strconv.Atoi(c.Param("cart_id"))
	productID, _ := strconv.Atoi(c.Param("item_id"))

	// get cart item by id from the database
	item, err := h.service.FindItem(uint(cartID), uint(productID))
	if err != nil {
		response := utils.APIResponse("Failed find product item from cart item", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// delete the cart item from the database
	err = h.service.DeleteItem(item)
	if err != nil {
		response := utils.APIResponse("failed to delete cart item", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// return a success message
	c.JSON(http.StatusOK, gin.H{"message": "cart item deleted"})
}
