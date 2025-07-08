package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaymentsByPin returns payment information by PIN (GET /api/v1/payments/:pin)
func GetPaymentsByPin(c *gin.Context) {
	pin := c.Param("pin")
	if pin == "" {
		auth.ErrorResponse(c, http.StatusBadRequest, "PIN parameter is required")
		return
	}

	payments, err := repository.GetPaymentsByPin(pin)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payments})
}
