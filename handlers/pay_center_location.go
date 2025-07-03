package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPayCenterLocationsHandler(c *gin.Context) {
	regno := c.Query("regno")
	payCenterIDStr := c.Query("pay_center_id")
	db := database.DB
	if regno != "" {
		data, err := repository.GetPayCenterLocationsByRegno(db, regno)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}
	if payCenterIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pay_center_id or regno is required"})
		return
	}
	payCenterID, err := strconv.ParseInt(payCenterIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pay_center_id"})
		return
	}
	data, err := repository.GetPayCenterLocationsByPayCenterID(db, payCenterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
