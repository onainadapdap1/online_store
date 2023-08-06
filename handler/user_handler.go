package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/helpers"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type UserHandlerInterface interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
	GetUserByID(c *gin.Context)
}

type userHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) UserHandlerInterface {
	return &userHandler{service: service}
}

// Register User godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param input dtos.RegisterUserInput body dtos.RegisterUserInput{} true "register user"
// @Success 200 {object} dtos.UserRegisterFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/v1/register [post]
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


// Login User godoc
// @Summary Login user
// @Description User login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param input dtos.LoginUserInput body dtos.LoginUserInput{} true "Login user input"
// @Success 200 {object} dtos.UserRegisterFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/v1/login [post]
func (h *userHandler) LoginUser(c *gin.Context) {
	var input dtos.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := utils.APIResponse("Login failed input user", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	loggedInUser, err := h.service.LoginUser(input)
	fmt.Println("error : ", err)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := utils.APIResponse("Login failed credential", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := helpers.GenerateToken(loggedInUser.ID, loggedInUser.Email)
	if err != nil {
		response := utils.APIResponse("Login failed when generate token", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := dtos.FormatUserLogin(loggedInUser, token)
	response := utils.APIResponse("Successfully loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}


// Fetch user login godoc
// @Summary Fetch user login
// @Description Fetch user login
// @Tags users
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/users/fetch [get]
func (h *userHandler) GetUserByID(c *gin.Context) {
	userData := c.MustGet("currentuser").(models.User)
	userID := userData.ID

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response := utils.APIResponse("failed to fetch user login", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := dtos.FormatUserLogin(user, "")

	response := utils.APIResponse("Successfully fetch user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}