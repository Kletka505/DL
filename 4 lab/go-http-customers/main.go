package main

import (
	"go-http-customers/database"
	"go-http-customers/handlers"
	"go-http-customers/middlewares"
	"go-http-customers/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	database.DB.AutoMigrate(&models.Customer{})

	r := gin.Default()

	r.POST("/customers", handlers.CreateCustomer)

	r.POST("/login", handlers.LoginCustomer)

	// Авторизованные маршруты
	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	authorized.GET("/customers/:id", handlers.GetCustomerData)

	r.Run(":8080")
}
