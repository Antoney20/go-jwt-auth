package controller

import (
	"net/http"
	"strings"

	"example.com/jwt-auth/tokens"
	"github.com/gin-gonic/gin"
)


func RefreshToken(c *gin.Context) {
	var refreshToken string
	refreshToken = c.GetHeader("Refresh-Token")

	if refreshToken == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) == 2 {
				refreshToken = splitToken[1]
			}
		}
	}
	if refreshToken == "" {
		refreshToken = c.Query("refresh_token")
	}


	if refreshToken == "" {
		refreshToken = c.PostForm("refresh_token")
	}

	if refreshToken == "" {
		var jsonInput struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&jsonInput); err == nil {
			refreshToken = jsonInput.RefreshToken
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	newAccessToken, err := tokens.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
