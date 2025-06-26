package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	ActiveValue int    `json:"activeValue"`
	MonitorId   string `json:"monitorId"`
	Code        string `json:"code"`
	FullName    string `json:"fullName"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	externalURL := "https://st-tais.mta.mn/rest/tais-ims-service/token/login"
	externalReq, err := http.NewRequest("POST", externalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create external request"})
		return
	}
	externalReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	externalResp, err := client.Do(externalReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to external API"})
		return
	}
	defer externalResp.Body.Close()

	body, err := io.ReadAll(externalResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read external response"})
		return
	}

	if externalResp.StatusCode != http.StatusOK {
		c.JSON(externalResp.StatusCode, gin.H{"error": string(body)})
		return
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse external response"})
		return
	}

	// Also try to extract code and fullName if not already set
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err == nil {
		if v, ok := raw["code"].(string); ok {
			loginResp.Code = v
		}
		if v, ok := raw["fullName"].(string); ok {
			loginResp.FullName = v
		}
	}

	c.JSON(http.StatusOK, loginResp)
}
