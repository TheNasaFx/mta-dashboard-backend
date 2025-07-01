package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ErrorResponse returns a unified error structure
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// AuthMiddleware validates JWT (without signature), sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ErrorResponse(c, http.StatusUnauthorized, "Authorization header missing or invalid")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}
		c.Set("user", claims)
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
