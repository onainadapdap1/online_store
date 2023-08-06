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

type PaymentMethodHandlerInterface interface {
	CreatePaymentMethod(c *gin.Context)
}

type paymentMethodHandler struct {
	service service.PaymentMethodServiceInterface
	paymentCategoryService service.PaymentCategoryServiceInterface
}

func NewPaymentMethodHandler(service service.PaymentMethodServiceInterface, paymentCategoryService service.PaymentCategoryServiceInterface) PaymentMethodHandlerInterface {
	return &paymentMethodHandler{
		service: service,
		paymentCategoryService: paymentCategoryService,
	}
}


// Create payment method godoc
// @Summary Create payment method
// @Description Create a new payment method with given data category payment id, method name, owner name and and number
// @Tags paymentmethods
// @Produce json
// @Param category_payment_id formData int true "category payment id of the payment method"
// @Param method_name formData string true "method name of the payment method"
// @Param owner_name formData string true "owner name of the payment method"
// @Param number formData int true "number of the payment method"
// @Success 200 {object} dtos.PaymentMethodFormatter
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/paymentmethods [post]
func (h *paymentMethodHandler) CreatePaymentMethod(c *gin.Context) {
	categoryPaymentID, err := strconv.Atoi(c.PostForm("category_payment_id"))
	if err != nil {
		response := utils.APIResponse("failed to convert category payment id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	methodName := c.PostForm("method_name")
	ownerName := c.PostForm("owner_name")
	number := c.PostForm("number")
	

	currentUser := c.MustGet("currentUser").(models.User)
	// userID := currentUser.ID

	paymentCategory, err := h.paymentCategoryService.GetPaymentCategoryByID(uint(categoryPaymentID))
	if err != nil {
		response := utils.APIResponse("failed to get payment category detail", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createPaymentMethod := dtos.CreatePaymentMethodInput{
		CategoryPaymentID: uint(categoryPaymentID),
		MethodName: methodName,
		OwnerName: ownerName,
		Number: number,
		User: currentUser,
		PaymentCategory: paymentCategory,
	}

	newPaymentMethod, err := h.service.CreatePaymentMethod(createPaymentMethod)
	if err != nil {
		log.Printf("failed to create payment method : %v", err)
		response := utils.APIResponse("failed to create payment method", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to create payment method", http.StatusCreated, "success", dtos.FormateCreatePaymentMethod(newPaymentMethod))
	c.JSON(http.StatusCreated, response)
}