package handlers

import (
	"dashboard-backend/database"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// safeParseFloat converts a string to float64, handling common formats
func safeParseFloat(s string) float64 {
	if s == "" {
		return 0
	}

	// Clean the string - remove extra spaces and replace comma with dot
	cleaned := strings.TrimSpace(s)
	cleaned = strings.Replace(cleaned, ",", ".", -1)

	// Check if it's a valid number format
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	if !re.MatchString(cleaned) {
		return 0
	}

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return val
}

// GetRealEstateStatistics returns real estate statistics
func GetRealEstateStatistics(c *gin.Context) {
	var totalCount int

	// 1. Нийт үл хөдлөхийн тоо - PAY_MARKET MRCH_REGNO -> V_TPI_PROPERTY_XYP_DATA_OWNER REG_NUM
	countQuery := `
		SELECT COUNT(*)
		FROM GPS.PAY_MARKET pm
		INNER JOIN GPS.V_TPI_PROPERTY_XYP_DATA_OWNER pxy ON pm.MRCH_REGNO = pxy.REG_NUM
	`
	err := database.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting real estate: " + err.Error()})
		return
	}

	// 2. Get all PROPERTY_SIZE and PROPERTY_VALUE as strings for processing in Go
	dataQuery := `
		SELECT pxy.PROPERTY_SIZE, pxy.PROPERTY_VALUE
		FROM GPS.PAY_MARKET pm
		INNER JOIN GPS.V_TPI_PROPERTY_XYP_DATA_OWNER pxy ON pm.MRCH_REGNO = pxy.REG_NUM
		WHERE pxy.PROPERTY_SIZE IS NOT NULL AND pxy.PROPERTY_VALUE IS NOT NULL
	`
	rows, err := database.DB.Query(dataQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying property data: " + err.Error()})
		return
	}
	defer rows.Close()

	var totalSize, totalValue float64
	for rows.Next() {
		var sizeStr, valueStr *string
		err := rows.Scan(&sizeStr, &valueStr)
		if err != nil {
			continue // Skip problematic rows
		}

		if sizeStr != nil {
			totalSize += safeParseFloat(*sizeStr)
		}
		if valueStr != nil {
			totalValue += safeParseFloat(*valueStr)
		}
	}

	// 4. Дундаж мвк ийн үнэлгээ (PROPERTY_VALUE / PROPERTY_SIZE)
	var avgPricePerSqm float64
	if totalSize > 0 {
		avgPricePerSqm = totalValue / totalSize
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_count":       totalCount,
			"total_size":        totalSize,
			"total_value":       totalValue,
			"avg_price_per_sqm": avgPricePerSqm,
		},
	})
}

// GetRealEstateData returns paginated real estate data with search
func GetRealEstateData(c *gin.Context) {
	// Pagination parameters
	page := 1
	limit := 50 // Default page size
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	offset := (page - 1) * limit

	// Build WHERE clause for search
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = "WHERE UPPER(pxy.REG_NUM) LIKE UPPER(:1)"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Count total records
	countQuery := `
		SELECT COUNT(*)
		FROM GPS.PAY_MARKET pm
		INNER JOIN GPS.V_TPI_PROPERTY_XYP_DATA_OWNER pxy ON pm.MRCH_REGNO = pxy.REG_NUM
		` + whereClause

	var totalRecords int
	err := database.DB.QueryRow(countQuery, args...).Scan(&totalRecords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting records: " + err.Error()})
		return
	}

	// Get paginated data - Return strings for numeric fields
	dataQuery := `
		SELECT 
			pxy.REG_NUM,
			pxy.LAST_NAME,
			pxy.FIRST_NAME,
			pxy.PROPERTY_NUMBER,
			pxy.PROPERTY_SIZE,
			pxy.PROPERTY_VALUE,
			pxy.FULL_ADDRESS,
			pxy.PROPERTY_TYPE
		FROM GPS.PAY_MARKET pm
		INNER JOIN GPS.V_TPI_PROPERTY_XYP_DATA_OWNER pxy ON pm.MRCH_REGNO = pxy.REG_NUM
		` + whereClause + `
		ORDER BY pxy.REG_NUM
		OFFSET :` + strconv.Itoa(argIndex) + ` ROWS FETCH NEXT :` + strconv.Itoa(argIndex+1) + ` ROWS ONLY
	`

	args = append(args, offset, limit)
	rows, err := database.DB.Query(dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var regNum, lastName, firstName, propertyNumber, fullAddress, propertyType *string
		var propertySizeStr, propertyValueStr *string

		err := rows.Scan(&regNum, &lastName, &firstName, &propertyNumber, &propertySizeStr, &propertyValueStr, &fullAddress, &propertyType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		// Convert string values to float64 safely
		var propertySize, propertyValue float64
		if propertySizeStr != nil {
			propertySize = safeParseFloat(*propertySizeStr)
		}
		if propertyValueStr != nil {
			propertyValue = safeParseFloat(*propertyValueStr)
		}

		// Combine first and last name
		ownerName := ""
		if lastName != nil && firstName != nil {
			ownerName = strings.TrimSpace(*lastName + " " + *firstName)
		} else if lastName != nil {
			ownerName = *lastName
		} else if firstName != nil {
			ownerName = *firstName
		}

		result := map[string]interface{}{
			"reg_num":         regNum,
			"owner_name":      ownerName,
			"property_number": propertyNumber,
			"property_size":   propertySize,
			"property_value":  propertyValue,
			"full_address":    fullAddress,
			"property_type":   propertyType,
		}
		results = append(results, result)
	}

	totalPages := (totalRecords + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"records":       results,
			"total_records": totalRecords,
			"total_pages":   totalPages,
			"current_page":  page,
			"page_size":     limit,
		},
	})
}
