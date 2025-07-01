package handlers

import (
	"dashboard-backend/auth"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProxyWorkerProfileHandler proxies the worker profile request to the external API with improved validation and error handling
func ProxyWorkerProfileHandler(c *gin.Context) {
	workerCode := c.Query("workerCode")
	isPrimary := c.DefaultQuery("isPrimary", "1")
	if workerCode == "" {
		auth.ErrorResponse(c, http.StatusBadRequest, "workerCode is required")
		return
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		auth.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
		return
	}

	externalURL := "https://st-tais.mta.mn/rest/tais-hrm-service/sql/workerPositionList/get?workerCode=" + workerCode + "&isPrimary=" + isPrimary

	req, err := http.NewRequest("GET", externalURL, nil)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to create request")
		return
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadGateway, "Failed to connect to external API")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to read response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		auth.ErrorResponse(c, resp.StatusCode, string(body))
		return
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(body, &profile); err == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}
