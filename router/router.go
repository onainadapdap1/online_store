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

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	api := router.Group("/api/v1")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	
	userRouter := api.Group("/users")
	{
		userRouter.Use(middlewares.Authentication())
		userRouter.GET("fetch", userAuthorization(userService), userHandler.GetUserByID)
	}
	categoryRouter := api.Group("/categories")
	{
		categoryRouter.GET("", categoryHandler.FindAllCategory)
		categoryRouter.GET("category/:slug", categoryHandler.FindBySlug)
		categoryRouter.Use(middlewares.Authentication())
		categoryRouter.POST("category", userAdminAuthorization(userService), categoryHandler.CreateCategory)
		categoryRouter.PUT("category/:slug", userAdminAuthorization(userService), categoryHandler.UpdateCategory)
		categoryRouter.DELETE("category/:id", userAdminAuthorization(userService), categoryHandler.DeleteCategoryByID)
	}
	productRouter := api.Group("/products")
	{
		productRouter.GET("", productHandler.FindAllProduct)
		productRouter.GET("product/:slug", productHandler.FindProductBySlug)
		productRouter.Use(middlewares.Authentication())
		productRouter.POST("product", userAdminAuthorization(userService), productHandler.CreateProduct)
		productRouter.PUT("product/:slug", userAdminAuthorization(userService), productHandler.UpdateProduct)
	}

	return router
}

func userAuthorization(userService service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["user_id"].(float64))

		user, err := userService.GetUserByID(userId)
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


func userAdminAuthorization(userService service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["user_id"].(float64))

		user, err := userService.GetUserByID(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": "data doesn't exist",
			})
		}

		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "you are now allowed to access this data",
			})
			return
		}
		c.Set("currentUser", user)
		c.Next()
	}
}