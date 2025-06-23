package auth

import (
	"github.com/gin-gonic/gin"
)

// TokenMiddleware нь authentication хийхгүй, шууд дараагийн handler руу дамжуулна.
func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
