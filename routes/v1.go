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
		v1.GET("/buildings/:id", handlers.GetBuildingByID)

		// Building floors and organizations endpoints
		v1.GET("/buildings/:id/floors", handlers.GetFloors)
		v1.GET("/buildings/:id/floors/:floor/organizations", handlers.GetOrganizations)
		v1.GET("/buildings/:id/organizations", handlers.GetAllOrganizations)

		// Taxpayer summary endpoint

		// LandView endpoint
		v1.GET("/land-views", handlers.GetLandViewsHandler)

		// PropertyOwner endpoint
		v1.GET("/property-owners", handlers.GetPropertyOwnersHandler)

		// Payments endpoint
		v1.GET("/payments/:pin", handlers.GetPaymentsByPin)

		// AccountGeneralYear endpoint
		v1.GET("/account-general-years", handlers.GetAccountGeneralYearsHandler)

		// TubReportData endpoint (TIN query param)
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

		// Map data endpoints
		v1.GET("/map-data", handlers.GetMapDataHandler)
		v1.GET("/map-data-batch", handlers.GetMapDataBatchHandler)
		v1.GET("/pay-center-properties", handlers.GetPayCenterPropertiesHandler)

		// Markets endpoint
		v1.GET("/markets", handlers.GetMarketsByPayCenterID)

		// Debug endpoint for ebarimt mapping
		v1.GET("/ebarimt-debug", handlers.GetEbarimtDebug)

		// Ebarimt endpoint
		v1.GET("/ebarimt/:pin", handlers.GetEbarimtByPin)

		// Dashboard statistics endpoint
		v1.GET("/statistics", handlers.GetDashboardStatistics)

		// Registration statistics endpoint
		v1.GET("/registration-stats/:id", handlers.GetRegistrationStats)

		// Tax office statistics endpoint
		v1.GET("/tax-office-stats/:id", handlers.GetTaxOfficeStats)

		// Segment statistics endpoint
		v1.GET("/segment-stats/:id", handlers.GetSegmentStats)

		v1.GET("/organizations", handlers.ListOrganizations)

		// Organization detail endpoint
		v1.GET("/organization-detail/:regno", handlers.GetOrganizationDetail)
	}
}
