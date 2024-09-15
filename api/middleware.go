package api

import (
	"net/http"
	"strings"

	"file-sharing-system/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function for authenticating JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Split the header to get the token part
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// Parse the token
		token, err := utils.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if claims have the expected fields
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract userID from claims
		userID, ok := claims["id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store the user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
