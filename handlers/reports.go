package handlers

import (
	"dashboard-backend/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DistrictActivityReportData represents the structure for activity reports by district
type DistrictActivityReportData struct {
	OpTypeName string `json:"op_type_name"`
	DistCode   int    `json:"dist_code"`
	DistName   string `json:"dist_name"`
	Count      int    `json:"count"`
}

// GetActivityReportsByDistrict returns activity reports grouped by district and activity type
func GetActivityReportsByDistrict(c *gin.Context) {
	// DIST_CODE болон DIST_NAME mapping
	districtNames := map[int]string{
		23: "Баганууд дүүрэг",
		24: "Багахангай дүүрэг",
		25: "Сүхбаатар дүүрэг",
		26: "Баянзүрх дүүрэг",
		27: "Налайх дүүрэг",
		28: "Сонгинохайрхан дүүрэг",
		29: "Чингэлтэй дүүрэг",
		34: "Хан-Уул дүүрэг",
		35: "Баянгол дүүрэг",
	}

	query := `
		SELECT 
			OP_TYPE_NAME,
			DIST_CODE,
			COUNT(ID) as cnt
		FROM GPS.PAY_MARKET 
		WHERE DIST_CODE IN (23,24,25,26,27,28,29,34,35)
			AND STATUS = 1
			AND OP_TYPE_NAME IS NOT NULL
		GROUP BY OP_TYPE_NAME, DIST_CODE
		ORDER BY OP_TYPE_NAME, DIST_CODE
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database query error: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var results []DistrictActivityReportData
	for rows.Next() {
		var opTypeName sql.NullString
		var distCode int
		var count int

		if err := rows.Scan(&opTypeName, &distCode, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Scan error: " + err.Error(),
			})
			return
		}

		// DIST_NAME mapping хийх
		distName, exists := districtNames[distCode]
		if !exists {
			distName = "Тодорхойгүй дүүрэг"
		}

		// OP_TYPE_NAME null эсэхийг шалгах
		if opTypeName.Valid && opTypeName.String != "" {
			result := DistrictActivityReportData{
				OpTypeName: opTypeName.String,
				DistCode:   distCode,
				DistName:   distName,
				Count:      count,
			}
			results = append(results, result)
		}
	}

	// Total count нэмж буцаах
	var totalCount int
	totalQuery := `
		SELECT SUM(cnt) AS total_count
		FROM (
			SELECT COUNT(id) AS cnt
			FROM GPS.PAY_MARKET a
			WHERE DIST_CODE IN (23,24,25,26,27,28,29,34,35)
				AND STATUS = 1
				AND OP_TYPE_NAME IS NOT NULL
			GROUP BY OP_TYPE_NAME
		) sub
	`

	err = database.DB.QueryRow(totalQuery).Scan(&totalCount)
	if err != nil {
		totalCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        results,
		"total_count": totalCount,
	})
}

// GetActivityReportsByPayCenter returns activity reports for a specific PAY_CENTER_ID
func GetActivityReportsByPayCenter(c *gin.Context) {
	payCenterId := c.Param("id")
	if payCenterId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "PAY_CENTER_ID required",
		})
		return
	}

	query := `
		SELECT 
			OP_TYPE_NAME,
			COUNT(ID) as cnt
		FROM GPS.PAY_MARKET 
		WHERE PAY_CENTER_ID = :1
			AND STATUS = 1
			AND OP_TYPE_NAME IS NOT NULL
		GROUP BY OP_TYPE_NAME
		ORDER BY OP_TYPE_NAME
	`

	rows, err := database.DB.Query(query, payCenterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database query error: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	type ActivityReportData struct {
		OpTypeName string `json:"op_type_name"`
		Count      int    `json:"count"`
	}

	var results []ActivityReportData
	for rows.Next() {
		var opTypeName sql.NullString
		var count int

		if err := rows.Scan(&opTypeName, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Scan error: " + err.Error(),
			})
			return
		}

		// OP_TYPE_NAME null эсэхийг шалгах
		if opTypeName.Valid && opTypeName.String != "" {
			result := ActivityReportData{
				OpTypeName: opTypeName.String,
				Count:      count,
			}
			results = append(results, result)
		}
	}

	// Total count нэмж буцаах
	var totalCount int
	totalQuery := `
		SELECT COUNT(ID) AS total_count
		FROM GPS.PAY_MARKET 
		WHERE PAY_CENTER_ID = :1
			AND STATUS = 1
			AND OP_TYPE_NAME IS NOT NULL
	`

	err = database.DB.QueryRow(totalQuery, payCenterId).Scan(&totalCount)
	if err != nil {
		totalCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"data":          results,
		"total_count":   totalCount,
		"pay_center_id": payCenterId,
	})
}
