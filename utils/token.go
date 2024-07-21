// utils/token.go

package utils

import (
	"library/task/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key")

// Claims represents the JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token
func GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken verifies the JWT token
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
