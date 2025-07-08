package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTubReportDataHandler(c *gin.Context) {
	db := database.DB
	tin := c.Query("tin")
	data, err := repository.GetTubReportDataByTIN(db, tin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
