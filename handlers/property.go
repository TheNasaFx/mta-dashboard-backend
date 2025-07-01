package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListPropertiesByOrg(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID")
		return
	}
	properties, err := repository.GetPropertiesByPayCenterID(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert sql.Null* fields to null in JSON if invalid
	var result []map[string]interface{}
	for _, p := range properties {
		item := map[string]interface{}{
			"id":            p.ID,
			"pay_center_id": p.PayCenterID,
			"updated_date":  nil,
			"property_type": nil,
			"owner_regno":   nil,
			"property_size": nil,
			"rent_amount":   nil,
		}
		if p.UpdatedDate.Valid {
			item["updated_date"] = p.UpdatedDate.String
		}
		if p.PropertyType.Valid {
			item["property_type"] = p.PropertyType.String
		}
		if p.OwnerRegno.Valid {
			item["owner_regno"] = p.OwnerRegno.String
		}
		if p.PropertySize.Valid {
			item["property_size"] = p.PropertySize.Float64
		}
		if p.RentAmount.Valid {
			item["rent_amount"] = p.RentAmount.Float64
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// CreateProperty creates a new property (POST /api/v1/property)
func CreateProperty(c *gin.Context) {
	var input repository.PropertyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	created, err := repository.CreateProperty(input)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": created})
}

// UpdateProperty updates a property (PUT /api/v1/property/:id)
func UpdateProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var input repository.PropertyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated, err := repository.UpdateProperty(uint(id), input)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// DeleteProperty deletes a property (DELETE /api/v1/property/:id)
func DeleteProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = repository.DeleteProperty(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}
