package handlers

import (
	"dashboard-backend/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPropertyOwnersHandler(c *gin.Context) {
	regNum := c.Query("reg_num")
	fmt.Printf("Property owner query for reg_num: %s\n", regNum)
	if regNum == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reg_num parameter is required"})
		return
	}

	query := `
		SELECT PROPERTY_NUMBER, FULL_ADDRESS, PROPERTY_TYPE, TO_NUMBER(REPLACE(PROPERTY_SIZE, ',', '.')) AS PROPERTY_SIZE, CREATED_DATE, REG_NUM
		FROM GPS.V_TPI_PROPERTY_XYP_DATA_OWNER 
		WHERE REG_NUM = :1`
	fmt.Printf("Executing query: %s with reg_num: %s\n", query, regNum)
	rows, err := database.DB.Query(query, regNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var propertyNumber, fullAddress, propertyType, createdDate, regNumResult *string
		var propertySize *float64

		if err := rows.Scan(&propertyNumber, &fullAddress, &propertyType, &propertySize, &createdDate, &regNumResult); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"property_number": propertyNumber,
			"full_address":    fullAddress,
			"property_type":   propertyType,
			"property_size":   propertySize,
			"created_date":    createdDate,
			"reg_num":         regNumResult,
		}
		results = append(results, result)
	}

	fmt.Printf("Found %d properties for reg_num: %s\n", len(results), regNum)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
