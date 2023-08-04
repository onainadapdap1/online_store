package router

import (
	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/driver"
	"github.com/onainadapdap1/online_store/handler"
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

	return router
}