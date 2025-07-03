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

		// Taxpayer summary endpoint

		// LandView endpoint
		v1.GET("/land-views", handlers.GetLandViewsHandler)

		// PropertyOwner endpoint
		v1.GET("/property-owners", handlers.GetPropertyOwnersHandler)

		// Payments endpoint
		v1.GET("/payments", handlers.GetPaymentsHandler)

		// AccountGeneralYear endpoint
		v1.GET("/account-general-years", handlers.GetAccountGeneralYearsHandler)

		// TubReportData endpoint
		v1.GET("/tub-report-data", handlers.GetTubReportDataHandler)

		// TaxAuditPaper endpoint
		v1.GET("/tax-audit-papers", handlers.GetTaxAuditPapersHandler)

		// TaxAuditViolation endpoint
		v1.GET("/tax-audit-violations", handlers.GetTaxAuditViolationsHandler)

		// TaxAuditPenalty endpoint
		v1.GET("/tax-audit-penalties", handlers.GetTaxAuditPenaltiesHandler)

		// TubHrmWorkPerformance endpoint
		v1.GET("/tub-hrm-work-performance", handlers.GetTubHrmWorkPerformancesHandler)

		// TubAuditQrHistory endpoint
		v1.GET("/tub-audit-qr-history", handlers.GetTubAuditQrHistoriesHandler)

		// PayCenterLocation endpoint
		v1.GET("/pay-center-locations", handlers.GetPayCenterLocationsHandler)

		v1.GET("/organizations", handlers.ListOrganizations)
	}
}
