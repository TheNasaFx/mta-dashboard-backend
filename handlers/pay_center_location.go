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
	grouped := c.Query("grouped")
	db := database.DB

	// Check if grouped is requested
	if grouped == "true" {
		data, err := repository.GetPayCenterLocationsGrouped(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
		return
	}

	if regno != "" {
		data, err := repository.GetPayCenterLocationsByRegno(db, regno)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
		return
	}
	if payCenterIDStr != "" {
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
		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
		return
	}

	// Get all pay center locations
	data, err := repository.GetAllPayCenterLocations(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
