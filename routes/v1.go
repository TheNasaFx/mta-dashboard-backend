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
			organizationGroup.GET(":id/property", handlers.ListPropertiesByOrg)
			organizationGroup.POST("", handlers.CreateOrganization)
			organizationGroup.PUT(":id", handlers.UpdateOrganization)
			organizationGroup.DELETE(":id", handlers.DeleteOrganization)
		}

		propertyGroup := v1.Group("/property", auth.TokenMiddleware())
		{
			propertyGroup.POST("", handlers.CreateProperty)
			propertyGroup.PUT(":id", handlers.UpdateProperty)
			propertyGroup.DELETE(":id", handlers.DeleteProperty)
		}

		v1.POST("/login", handlers.LoginHandler)
		v1.GET("/organization/:id/market", auth.TokenMiddleware(), handlers.ListMarketsByOrg)

		// Proxy endpoint for external profile API
		v1.GET("/proxy/worker-profile", auth.TokenMiddleware(), handlers.ProxyWorkerProfileHandler)

		// --- Add these for centers and buildings ---
		v1.GET("/centers", handlers.GetCenters)
		v1.GET("/buildings", handlers.GetBuildings)
	}
}
