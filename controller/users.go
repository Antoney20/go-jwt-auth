package controller

import (
	"log"
	"net/http"

	"example.com/jwt-auth/config"
	"example.com/jwt-auth/models"
	"example.com/jwt-auth/tokens"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data now"})
		return
	}

	if err := model.ValidatePhoneNumber(user.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser model.User
	if err := config.DB.Where("phone_number = ?", user.PhoneNumber).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is already registered"})
		return
	}

	if err := user.Validate(config.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.HashPassword()

	// save
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var inputdata struct {
		Identifier string `json:"identifier"` 
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&inputdata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var user model.User

	log.Printf("Login attempt for Identifier: %s", inputdata.Identifier)

	if err := config.DB.Where("username = ? OR email = ?", inputdata.Identifier, inputdata.Identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password provided"})
		return
	}

	if !user.CheckPassword(inputdata.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := tokens.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful!",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
