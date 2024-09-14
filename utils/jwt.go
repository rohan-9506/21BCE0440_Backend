package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT_SECRET_KEY should be defined in your .env or config file
var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

// GenerateJWT generates a new JWT token for the given user ID
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT parses and validates the JWT token
func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(JWT_SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Custom error for invalid token
var ErrInvalidToken = jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
