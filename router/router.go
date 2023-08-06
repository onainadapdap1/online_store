package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/handler"
	"github.com/onainadapdap1/online_store/middlewares"
	"github.com/onainadapdap1/online_store/repository"
	"github.com/onainadapdap1/online_store/service"
	"github.com/gin-contrib/cors"

	_ "github.com/onainadapdap1/online_store/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 			Online Store API
// @version         1.0
// @description     This is service for my Online Store API assignment.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     https://onainadapdap1.github.io/
// @contact.email   nadapdaponai21@gmail.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func Router() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://localhost:8080"},
        AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return origin == "https://github.com"
        },
        MaxAge: 12 * time.Hour,
    }))
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

	paymentCategoryRepository := repository.NewPaymentCategoryRepository(db)
	paymentCategoryService := service.NewPaymentCategoryService(paymentCategoryRepository)
	paymentCategoryHandler := handler.NewPaymentCategoryHandler(paymentCategoryService)

	paymentMethodRepository := repository.NewPaymentMethodRepository(db)
	paymentMethodService := service.NewPaymentMethodService(paymentMethodRepository)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodService, paymentCategoryService)

	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepository)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository)
	orderHandler := handler.NewOrderHandler(orderService)

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
	paymentCategoryRouter := api.Group("/paymentcategories")
	{
		paymentCategoryRouter.GET("paymentcategory/:slug", paymentCategoryHandler.GetPaymentCategoryBySlug)
		paymentCategoryRouter.Use(middlewares.Authentication())
		paymentCategoryRouter.POST("", userAdminAuthorization(userService), paymentCategoryHandler.CreatePaymentCategory)
	}
	paymentMethodRouter := api.Group("/paymentmethods")
	{
		paymentMethodRouter.Use(middlewares.Authentication())
		paymentMethodRouter.POST("", userAdminAuthorization(userService), paymentMethodHandler.CreatePaymentMethod)
	}
	cartRouter := api.Group("/carts")
	{
		cartRouter.Use(middlewares.Authentication())
		cartRouter.POST("/cart", userAuthorization(userService), cartHandler.AddProductToCart)
		cartRouter.PUT("/cart/:cart_id/productID/:item_id", userAuthorization(userService), cartHandler.UpdateCartItemQuantity)
		cartRouter.DELETE("/cart/:cart_id/productID/:item_id", userAuthorization(userService), cartHandler.DeleteCartItem)
		cartRouter.GET("", userAuthorization(userService), cartHandler.GetAllUserCartItems)
	}
	orderRouter := api.Group("/orders")
	{
		orderRouter.Use(middlewares.Authentication())
		orderRouter.POST("/", userAuthorization(userService), orderHandler.Checkout)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
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