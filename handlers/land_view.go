package handlers

import (
	"dashboard-backend/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLandViewsHandler(c *gin.Context) {
	pin := c.Query("pin")
	if pin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pin parameter is required"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT CERTIFICATE_NO, AU1_NAME, AU2_NAME, AU3_NAME, ADDRESS_STREETNAME, 
		       DECISION_NO, AREA_M2, ADDRESS_KHASHAA, PIN 
		FROM GPS.V_E_TUB_LAND_VIEW 
		WHERE PIN = :1`, pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var certificateNo, au1Name, au2Name, au3Name, addressStreetname *string
		var decisionNo, addressKhashaa, pinResult *string
		var areaM2 *float64

		if err := rows.Scan(&certificateNo, &au1Name, &au2Name, &au3Name, &addressStreetname,
			&decisionNo, &areaM2, &addressKhashaa, &pinResult); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		// Хаяг бүтээх
		var fullAddress string
		if au1Name != nil {
			fullAddress += *au1Name
		}
		if au2Name != nil {
			if fullAddress != "" {
				fullAddress += " "
			}
			fullAddress += *au2Name
		}
		if au3Name != nil {
			if fullAddress != "" {
				fullAddress += " "
			}
			fullAddress += *au3Name
		}
		if addressStreetname != nil {
			if fullAddress != "" {
				fullAddress += " "
			}
			fullAddress += *addressStreetname
		}

		// Гудамж (зөвхөн ADDRESS_STREETNAME ашиглана)
		var street string
		if addressStreetname != nil {
			street = *addressStreetname
		}

		result := map[string]interface{}{
			"certificate_no":  certificateNo,
			"full_address":    fullAddress,
			"decision_no":     decisionNo,
			"area_m2":         areaM2,
			"au1_name":        au1Name,
			"au2_name":        au2Name,
			"au3_name":        au3Name,
			"street":          street,
			"address_khashaa": addressKhashaa,
			"pin":             pinResult,
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
