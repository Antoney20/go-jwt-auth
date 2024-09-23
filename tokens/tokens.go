package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

// TokenClaims represents the custom claims for the JWT
type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateTokens generates access and refresh tokens
func GenerateTokens(userID uint) (string, string, error) {
	// Set token expiration times
	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	refreshTokenExpiration := time.Now().Add(7 * 24 * time.Hour)

	// Create access token
	accessTokenClaims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshTokenClaims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiration),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateToken checks if the token is valid
func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid or expired token")
}
