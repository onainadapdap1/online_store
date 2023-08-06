package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type PaymentCategoryHandlerInterface interface {
	CreatePaymentCategory(c *gin.Context)
	GetPaymentCategoryBySlug(c *gin.Context)
}

type paymentCategoryHandler struct {
	service service.PaymentCategoryServiceInterface
}

func NewPaymentCategoryHandler(service service.PaymentCategoryServiceInterface) PaymentCategoryHandlerInterface {
	return &paymentCategoryHandler{service: service}
}

// CreatePaymentCategory godoc
// @Summary Create Payment Category
// @Description Create Payment Category
// @Tags payment-categories
// @Accept json
// @Produce json
// @Param input body dtos.CreatePaymentCategoryInput{} true "Payment category input"
// @Success 200 {object} dtos.PaymentCategoryFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/paymentcategories [post]
func (h *paymentCategoryHandler) CreatePaymentCategory(c *gin.Context) {
	var input dtos.CreatePaymentCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := utils.APIResponse("Failed to load user input", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	createPaymentCategoryInput := dtos.CreatePaymentCategoryInput{
		CategoryName: input.CategoryName,
		User:         currentUser,
	}

	newPaymentCategory, err := h.service.CreatePaymentCategory(createPaymentCategoryInput)
	if err != nil {
		log.Printf("failed to create payment category: %v", err)
		// Check if the error is a duplicate key error
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			response := utils.APIResponse("Payment category already exists", http.StatusConflict, "error", nil)
			c.JSON(http.StatusConflict, response)
			return
		}

		response := utils.APIResponse("Failed to create payment category", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("Success to create payment category", http.StatusOK, "success", dtos.FormateCreatePaymentCategory(newPaymentCategory))
	c.JSON(http.StatusOK, response)
}

// Get payment category by slug godoc
// @Summary Get payment category by slug
// @Description Get payment category by slug
// @Tags payment-categories
// @Produce json
// @Param slug path string true "Get payment category by slug"
// @Success 200 {object} dtos.PaymentCategoryFormatter
// @Failure 400 {object} utils.Response
// @Router /api/v1/paymentcategories/paymentcategory/{slug} [get]
func (h *paymentCategoryHandler) GetPaymentCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// log.Println("slug : ", slug)
	paymentCategory, err := h.service.GetPaymentCategoryBySlug(slug)
	if err != nil {
		response := utils.APIResponse("failed to get detail payment category", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to get product detail", http.StatusOK, "success", dtos.FormateCreatePaymentCategory(paymentCategory))
	c.JSON(http.StatusOK, response)
}
