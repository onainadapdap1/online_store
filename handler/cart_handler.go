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


// Get All user item cart godoc
// @Summary Get All user item cart
// @Description Get All user item cart
// @Tags carts
// @Produce json
// @Success 200 {object} []dtos.CartItemFormatter{}
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/carts [get]
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


// Add product to cart godoc
// @Summary Add product to cart
// @Description Add product to cart
// @Tags carts
// @Accept json
// @Produce json
// @Param input body dtos.CreateCartItemInput{} true "add product to cart input"
// @Success 200 {object} dtos.CartItemFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/carts/cart [post]
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

// UpdateCartItemQuantity godoc
// @Summary Update cart item quantity
// @Description Update cart item quantity based on action (add/remove) and quantity
// @Tags carts
// @Accept json
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Param item_id path int true "Cart Item ID"
// @Param action formData string true "Action to perform (add/remove)"
// @Param quantity formData int true "Quantity to add/remove"
// @Success 200 {object} dtos.CartItemFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/carts/cart/{cart_id}/productID/{item_id} [put]
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


// Delete item from cart godoc
// @Summary Delete item from cart
// @Description Delete item from cart
// @Tags carts
// @Produce json
// @Param cart_id path int true "cart id"
// @Param item_id path int true "item id"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/carts/cart/{cart_id}/productID/{item_id} [delete]
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
