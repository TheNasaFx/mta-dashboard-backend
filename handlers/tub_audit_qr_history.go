package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTubAuditQrHistoriesHandler(c *gin.Context) {
	db := database.DB
	regno := c.Query("regno")
	if regno != "" {
		data, err := repository.GetTubAuditQrHistoriesByRegno(db, regno)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": len(data)})
		return
	}
	data, err := repository.GetTubAuditQrHistories(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
