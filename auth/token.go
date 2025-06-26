package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // Use the same key as token generation if you generate your own tokens

// AuthMiddleware checks only for the presence and format of the Bearer token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}
		// Do not parse or validate the token, just pass it through
		c.Next()
	}
}

// GetUserFromContext extracts user info from context (if set by AuthMiddleware)
func GetUserFromContext(c *gin.Context) (map[string]interface{}, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return nil, false
	}
	return claims, true
}
