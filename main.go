package main

import (
	"log"

	"example.com/jwt-auth/config"
	"example.com/jwt-auth/controller"
	"example.com/jwt-auth/middleware"
	"example.com/jwt-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// load env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Connect to the database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Auto migrate models
	if err := config.DB.AutoMigrate(&model.User{}, &model.Profile{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Migrations successful")

	router := gin.Default()

	router.POST("/register", controller.RegisterUser)
	router.POST("/login", controller.LoginUser)
	router.POST("/refresh-token", controller.RefreshToken)
	router.GET("/profile", middleware.AuthenticateMiddleware(), controller.GetProfile)
	router.POST("/profile", middleware.AuthenticateMiddleware(), controller.CreateProfile)
	router.PUT("/profile", middleware.AuthenticateMiddleware(), controller.UpdateProfile)


	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
