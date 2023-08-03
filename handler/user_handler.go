package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type UserHandlerInterface interface {
	RegisterUser(c *gin.Context)
}

type userHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) UserHandlerInterface {
	return &userHandler{service: service}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input dtos.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.service.RegisterUser(input)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userFormatReg := dtos.FormatUserRegister(&newUser)
	response := utils.APIResponse("Account has been registered", http.StatusCreated, "success", userFormatReg)
	c.JSON(http.StatusCreated, response)

}
