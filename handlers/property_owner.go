package handlers

import (
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPropertyOwnersHandler(c *gin.Context) {
	regno := c.Query("regno")
	entId := c.Query("ent_id")
	db := database.DB
	data, err := repository.GetPropertyOwners(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entId != "" {
		var filtered []interface{}
		for _, owner := range data {
			if owner.ENT_ID.Valid && strconv.FormatInt(owner.ENT_ID.Int64, 10) == entId {
				filtered = append(filtered, owner)
			}
		}
		c.JSON(http.StatusOK, filtered)
		return
	}
	if regno != "" {
		var filtered []interface{}
		for _, owner := range data {
			if owner.REG_NUM.Valid && owner.REG_NUM.String == regno {
				filtered = append(filtered, owner)
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
