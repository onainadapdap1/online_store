package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type OrderHandlerInterface interface {
	Checkout(c *gin.Context)
}

type orderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(service service.OrderServiceInterface) *orderHandler {
	return &orderHandler{service: service}
}

// Create order godoc
// @Summary Create order
// @Description Create order
// @Tags orders
// @Accept mpfd
// @Produce json
// @Param payment_category_id formData int true "payment category id"
// @Param payment_method_id formData int true "payment method id"
// @Param receiver_name formData string true "receiver name"
// @Param proof_of_payment formData file true "Image file of the proof of payment"
// @Success 200 {object} dtos.OrderFormatter
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/orders [post]
func (h *orderHandler) Checkout(c *gin.Context) {
	currentUser := c.MustGet("currentuser").(models.User)
	userID := currentUser.ID
	// userID := c.MustGet("userID")

	// Check if user has cart (UPDATE 1)
	cart, err := h.service.CheckIsUserHasACart(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has no cart"})
		return
	}

	// GET ALL CART ITEM (UPDATE 2)
	cartItems, err := h.service.GetAllCartItems(cart.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot get all item inside of cart item"})
		return
	}

	totalPrice := h.service.CountPriceInCartItems(cartItems)

	paymentCategoryID, _ := strconv.Atoi(c.PostForm("payment_category_id"))
	paymentMethodID, _ := strconv.Atoi(c.PostForm("payment_method_id"))
	receiverName := c.PostForm("receiver_name")
	proofOfPayment, err := c.FormFile("proof_of_payment")
	if err != nil {
		response := utils.APIResponse("failed to load image file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fileName := fmt.Sprintf("%d-%s", userID, proofOfPayment.Filename)

	dirPath := filepath.Join(".", "static", "images", "payments")
	filePath := filepath.Join(dirPath, fileName)
	// create directory if doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			response := utils.APIResponse("failed to upload product image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	// create file that will hold the image
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// open the temporary file that contains the uploaded image
	inputFile, err := proofOfPayment.Open()
	if err != nil {
		response := utils.APIResponse("failed to open product input image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	defer inputFile.Close()

	// copy the temporary image to the permanent location outputFile
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		response := utils.APIResponse("failed to copy product input image to permanent location", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// UPDATE
	paymentCID, err := h.service.GetPaymenCID(uint(paymentCategoryID))
	if err != nil {
		response := utils.APIResponse("failed to get payment category", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	paymentMID, err := h.service.GetPaymentMID(uint(paymentCategoryID))
	if err != nil {
		response := utils.APIResponse("failed to get payment method", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Buat Order baru dengan status "Pending"
	orderData := models.Order{
		UserID:            userID,
		PaymentCategoryID: uint(paymentCategoryID),
		PaymentMethodID:   uint(paymentMethodID),
		ReceiverName:      receiverName,
		ProofOfPayment:    filePath,
		TotalPrice:        totalPrice,
		Status:            "Pending",
		PaymentCategory:   paymentCID,
		PaymentMethod:     paymentMID,
		// OrderItems:        orderItems,
	}

	order, err := h.service.CreateOrder(orderData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}
	if len(cartItems) > 0 {
		status, err := h.service.CreateOrderItem(order, cartItems)
		if err != nil {
			c.JSON(status, gin.H{"error": "failed to create orderItem and update product"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "your cart is empty"})
		return
	}

	err = h.service.UpdateCartAndDeleteItems(cart.ID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update cart and delete cart item"})
		return
	}

	// Berikan feedback ke user dengan informasi tentang Order dan Payment
	response := utils.APIResponse("success to create order", http.StatusCreated, "success", dtos.FormatOrderDetail(order))
	c.JSON(http.StatusCreated, response)
}
