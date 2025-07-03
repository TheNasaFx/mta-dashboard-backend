package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaymentsHandler(c *gin.Context) {
	ownerID := c.Query("owner_id")
	db := database.DB
	data, err := repository.GetPayments(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ownerID != "" {
		var filtered []interface{}
		for _, item := range data {
			if item.OWNER_ID.Valid {
				if id, err := strconv.ParseInt(ownerID, 10, 64); err == nil && item.OWNER_ID.Int64 == id {
					filtered = append(filtered, item)
				}
			}
		}
		if len(filtered) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusOK, filtered)
		return
	}
	c.JSON(http.StatusOK, data)
}
