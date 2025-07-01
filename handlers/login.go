package handlers

import (
	"bytes"
	"dashboard-backend/auth"
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
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to marshal request")
		return
	}

	externalURL := "https://st-tais.mta.mn/rest/tais-ims-service/token/login"
	externalReq, err := http.NewRequest("POST", externalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to create external request")
		return
	}
	externalReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	externalResp, err := client.Do(externalReq)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadGateway, "Failed to connect to external API")
		return
	}
	defer externalResp.Body.Close()

	body, err := io.ReadAll(externalResp.Body)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to read external response")
		return
	}

	if externalResp.StatusCode != http.StatusOK {
		auth.ErrorResponse(c, externalResp.StatusCode, string(body))
		return
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, "Failed to parse external response")
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

	c.JSON(http.StatusOK, gin.H{"success": true, "data": loginResp})
}
