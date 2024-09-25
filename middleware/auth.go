package middleware

import (
	"errors"
	"net/http"
	"strings"

	"example.com/jwt-auth/tokens"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromToken(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("no authorization header")
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return 0, errors.New("invalid authorization header")
	}

	claims, err := tokens.ValidateToken(splitToken[1])
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}

func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort() 
			return
		}

		c.Set("userID", userID)
		c.Next() 
	}
}
