package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLandViewsHandler(c *gin.Context) {
	pin := c.Query("pin")
	db := database.DB
	data, err := repository.GetLandViews(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pin != "" {
		var filtered []interface{}
		for _, lv := range data {
			if lv.PIN.Valid && lv.PIN.String == pin {
				filtered = append(filtered, lv)
			}
		}
		c.JSON(http.StatusOK, filtered)
		return
	}
	c.JSON(http.StatusOK, data)
}
