package tokens

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

type UserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenError struct {
	Message    string
	StatusCode int
}

func (e *TokenError) Error() string {
	return e.Message
}


func GenerateTokens(id uint, username string) (string, string, error) {
	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	refreshTokenExpiration := time.Now().Add(30 * 24 * time.Hour) // 30 days

	accessTokenClaims := UserClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	refreshTokenClaims := UserClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiration),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", errors.New("could not generate tokens")
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, &TokenError{"Invalid token", http.StatusUnauthorized}
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, &TokenError{"Invalid token", http.StatusUnauthorized}
}


func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Here you would typically check if the user still has permission to generate a new token
	// For example, check if the user still exists in the database and has the necessary permissions
	// For this example, we'll assume the check passes

	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	newClaims := UserClaims{
		ID:       claims.ID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	accessTokenString, err := newToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("could not generate new access token")
	}

	return accessTokenString, nil
}
// func RefreshAccessToken(refreshToken string) (string, error) {
// 	claims, err := ValidateToken(refreshToken)
// 	if err != nil {
// 		return "", errors.New("invalid refresh token")
// 	}

// 	// Here you would typically check if the user still has permission to generate a new token
// 	// For example, check if the user still exists in the database and has the necessary permissions
// 	// For this example, we'll assume the check passes

// 	accessTokenExpiration := time.Now().Add(15 * time.Minute)
// 	newClaims := UserClaims{
// 		ID:    claims.ID,
// 		First: claims.First,
// 		Last:  claims.Last,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
// 		},
// 	}

// 	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
// 	accessTokenString, err := newToken.SignedString(jwtSecretKey)
// 	if err != nil {
// 		return "", errors.New("could not generate new access token")
// 	}

// 	return accessTokenString, nil
// }
