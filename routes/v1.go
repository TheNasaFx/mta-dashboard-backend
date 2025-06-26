package routes

import (
	"dashboard-backend/auth"
	"dashboard-backend/handlers"
	"io"
	"net/http"

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

		v1.POST("/login", handlers.LoginHandler)
		v1.GET("/workerPositionList", auth.TokenMiddleware(), handlers.WorkerPositionListHandler)

		// Proxy endpoint for external profile API
		v1.GET("/proxy/worker-profile", func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			workerCode := c.Query("workerCode")
			if token == "" || workerCode == "" {
				c.JSON(400, gin.H{"error": "token or workerCode required"})
				return
			}
			client := &http.Client{}
			req, _ := http.NewRequest("GET", "https://st-tais.mta.mn/rest/tais-hrm-service/sql/workerPositionList/get?workerCode="+workerCode+"&isPrimary=1", nil)
			req.Header.Set("Authorization", token)
			resp, err := client.Do(req)
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to fetch"})
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		})
	}
}
