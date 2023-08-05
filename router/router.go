package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/handler"
	"github.com/onainadapdap1/online_store/middlewares"
	"github.com/onainadapdap1/online_store/repository"
	"github.com/onainadapdap1/online_store/service"
)

func Router() *gin.Engine {
	router := gin.Default()
	db, _ := driver.ConnectDB()
	
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	api := router.Group("/api/v1")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	
	userRouter := api.Group("/users")
	{
		userRouter.Use(middlewares.Authentication())
		userRouter.GET("fetch", userAuthorization(userService), userHandler.GetUserByID)
	}

	return router
}

func userAuthorization(userService service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := driver.GetDB()
		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["user_id"].(float64))
		// user := models.User{}

		user, err := userService.GetUserByID(userId)
		// err := db.Where("id = ?", userId).Find(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": "data doesn't exist",
			})
		}

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "User not found",
				"message": "Please try again",
			})
			return
		}
		c.Set("currentuser", user)
		c.Next()
	}
}