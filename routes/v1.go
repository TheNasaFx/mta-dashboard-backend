package routes

import (
	"dashboard-backend/auth"
	"dashboard-backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.Status(200)
		})

		organizationGroup := v1.Group("/organization", auth.TokenMiddleware())
		{
			organizationGroup.GET("", handlers.ListOrganization)
			organizationGroup.GET("/:id", handlers.GetOrganizationByID)
		}
	}
}
