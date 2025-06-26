package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WorkerPositionListHandler(c *gin.Context) {
	workerCode := c.Query("workerCode")
	isPrimary := c.Query("isPrimary")
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	externalURL := "https://st-tais.mta.mn/rest/tais-hrm-service/sql/workerPositionList/get?workerCode=" + workerCode + "&isPrimary=" + isPrimary

	req, err := http.NewRequest("GET", externalURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to external API"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	// Try to unmarshal as a single object
	var profile map[string]interface{}
	if err := json.Unmarshal(body, &profile); err == nil {
		c.JSON(http.StatusOK, gin.H{"profile": profile})
		return
	}

	// If not a single object, return as raw JSON
	c.Data(http.StatusOK, "application/json", body)
}
