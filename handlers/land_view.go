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

// GetAllLandDataHandler газрын мэдээллийн dashboard-д зориулсан функц
func GetAllLandDataHandler(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT PIN, AU2_NAME, AREA_M2, NAME 
		FROM GPS.V_E_TUB_LAND_VIEW`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var pin, au2Name, name *string
		var areaM2 *float64

		if err := rows.Scan(&pin, &au2Name, &areaM2, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"pii":      pin,
			"au2_name": au2Name,
			"area_m2":  areaM2,
			"name":     name,
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// GetLandPaymentDataHandler газрын татвар төлөлтийн мэдээлэл
func GetLandPaymentDataHandler(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT 
			l.PIN,
			l.AU2_NAME,
			p.BRANCH_CODE,
			p.BRANCH_NAME,
			p.SUB_BRANCH_CODE,
			p.SUB_BRANCH_NAME,
			p.TAX_TYPE_CODE,
			p.TAX_TYPE_NAME,
			p.AMOUNT
		FROM (SELECT PIN, AU2_NAME FROM GPS.V_E_TUB_LAND_VIEW WHERE ROWNUM <= 10000) l
		INNER JOIN GPS.V_E_TUB_PAYMENTS p ON l.PIN = p.PIN
		WHERE p.TAX_TYPE_CODE IN ('01030703', '01030731')
		AND ROWNUM <= 50000`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var pin, au2Name, branchCode, branchName *string
		var subBranchCode, subBranchName, taxTypeCode, taxTypeName *string
		var amount *float64

		if err := rows.Scan(&pin, &au2Name, &branchCode, &branchName,
			&subBranchCode, &subBranchName, &taxTypeCode, &taxTypeName, &amount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"pin":             pin,
			"au2_name":        au2Name,
			"branch_code":     branchCode,
			"branch_name":     branchName,
			"sub_branch_code": subBranchCode,
			"sub_branch_name": subBranchName,
			"tax_type_code":   taxTypeCode,
			"tax_type_name":   taxTypeName,
			"amount":          amount,
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
