package auth

import (
	"github.com/gin-gonic/gin"
)

// TokenMiddleware is now an alias for AuthMiddleware for real token validation
func TokenMiddleware() gin.HandlerFunc {
	return AuthMiddleware()
}
