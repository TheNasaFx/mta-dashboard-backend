package handlers

import (
	"dashboard-backend/database"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
)

// Cache structure for map data
type mapDataCache struct {
	data      map[string]interface{}
	timestamp time.Time
}

var (
	mapCache      = make(map[string]*mapDataCache)
	mapCacheMutex = sync.RWMutex{}
	cacheTimeout  = 2 * time.Minute // Reduced to 2 minutes for better responsiveness
)

func GetMapDataHandler(c *gin.Context) {
	payCenterID := c.Query("pay_center_id")
	if payCenterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pay_center_id parameter is required"})
		return
	}

	// Check cache first
	mapCacheMutex.RLock()
	if cached, exists := mapCache[payCenterID]; exists {
		if time.Since(cached.timestamp) < cacheTimeout {
			mapCacheMutex.RUnlock()
			c.JSON(http.StatusOK, gin.H{"success": true, "data": cached.data})
			return
		}
	}
	mapCacheMutex.RUnlock()

	// Convert to int for database query
	id, err := strconv.Atoi(payCenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay_center_id"})
		return
	}

	fmt.Printf("Getting map data for PAY_CENTER_ID: %d\n", id)

	// Debug: Check if PAY_CENTER exists
	var centerExists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER WHERE ID = :1", id).Scan(&centerExists)
	if err != nil {
		fmt.Printf("PAY_CENTER check error: %v\n", err)
	}
	fmt.Printf("PAY_CENTER exists: %d\n", centerExists)

	// Debug: Check total records in PAY_CENTER_PROPERTY
	var totalPropertyRecords int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER_PROPERTY").Scan(&totalPropertyRecords)
	if err == nil {
		fmt.Printf("Total PAY_CENTER_PROPERTY records: %d\n", totalPropertyRecords)
	}

	// Debug: Check records for this pay_center_id
	var propertyRecordsForCenter int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER_PROPERTY WHERE PAY_CENTER_ID = :1", id).Scan(&propertyRecordsForCenter)
	if err == nil {
		fmt.Printf("PAY_CENTER_PROPERTY records for ID %d: %d\n", id, propertyRecordsForCenter)
	}

	// Debug: Check PAY_MARKET records
	var marketRecordsForCenter int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :1", id).Scan(&marketRecordsForCenter)
	if err == nil {
		fmt.Printf("PAY_MARKET records for ID %d: %d\n", id, marketRecordsForCenter)
	}

	// Query owner count
	var ownerCount int
	ownerQuery := `SELECT COUNT(DISTINCT OWNER_REGNO) FROM GPS.PAY_CENTER_PROPERTY WHERE PAY_CENTER_ID = :1`
	err = database.DB.QueryRow(ownerQuery, id).Scan(&ownerCount)
	if err != nil {
		fmt.Printf("Owner count query error: %v\n", err)
		ownerCount = 0
	}
	fmt.Printf("Owner count: %d\n", ownerCount)

	// Query activity operators (total records)
	var activityOperators int
	activityQuery := `SELECT COUNT(*) FROM GPS.PAY_CENTER_PROPERTY WHERE PAY_CENTER_ID = :1`
	err = database.DB.QueryRow(activityQuery, id).Scan(&activityOperators)
	if err != nil {
		fmt.Printf("Activity operators query error: %v\n", err)
		activityOperators = 0
	}
	fmt.Printf("Activity operators: %d\n", activityOperators)

	// Query total area
	var area float64
	areaQuery := `SELECT NVL(SUM(PROPERTY_SIZE), 0) FROM GPS.PAY_CENTER_PROPERTY WHERE PAY_CENTER_ID = :1`
	err = database.DB.QueryRow(areaQuery, id).Scan(&area)
	if err != nil {
		fmt.Printf("Area query error: %v\n", err)
		area = 0
	}
	fmt.Printf("Total area: %f\n", area)

	// Query land area from V_E_TUB_LAND_VIEW (AREA_M2)
	var landArea float64
	landAreaQuery := `
		SELECT NVL(SUM(
			NVL(AREA_M2, 0)
		), 0) as land_area
		FROM GPS.V_E_TUB_LAND_VIEW 
		WHERE PAY_CENTER_ID = :1
	`
	err = database.DB.QueryRow(landAreaQuery, id).Scan(&landArea)
	if err != nil {
		fmt.Printf("Land area query error: %v\n", err)
		// Try without GPS schema
		err2 := database.DB.QueryRow(`
			SELECT NVL(SUM(
				NVL(AREA_M2, 0)
			), 0) as land_area
			FROM V_E_TUB_LAND_VIEW 
			WHERE PAY_CENTER_ID = :1
		`, id).Scan(&landArea)
		if err2 != nil {
			fmt.Printf("Land area query error (alt): %v\n", err2)
			landArea = 0
		}
	}
	fmt.Printf("Total land area: %f\n", landArea)

	// Calculate unused area
	unusedArea := landArea - area
	fmt.Printf("Unused area: %f\n", unusedArea)

	// Query tenants count (PAY_MARKET records - unique tenants)
	var tenants int
	tenantsQuery := `SELECT COUNT(DISTINCT MRCH_REGNO) FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :1`
	err = database.DB.QueryRow(tenantsQuery, id).Scan(&tenants)
	if err != nil {
		fmt.Printf("Tenants query error: %v\n", err)
		tenants = 0
	}
	fmt.Printf("Tenants count: %d\n", tenants)

	result := map[string]interface{}{
		"owner_count":        ownerCount,
		"activity_operators": activityOperators,
		"area":               area,       // Ашиглагдаж байгаа талбай (мкв)
		"land_area":          landArea,   // Газрын талбай (мкв)
		"unused_area":        unusedArea, // Ашиглагдаагүй талбай (мкв)
		"tenants":            tenants,
	}

	// Update cache
	mapCacheMutex.Lock()
	mapCache[payCenterID] = &mapDataCache{
		data:      result,
		timestamp: time.Now(),
	}
	mapCacheMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetMapDataBatchHandler returns map data for multiple pay centers in one request
func GetMapDataBatchHandler(c *gin.Context) {
	payCenterIDsStr := c.Query("pay_center_ids")
	if payCenterIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pay_center_ids parameter is required"})
		return
	}

	// Parse comma-separated IDs
	idStrings := strings.Split(payCenterIDsStr, ",")
	var payCenterIDs []int
	for _, idStr := range idStrings {
		if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
			payCenterIDs = append(payCenterIDs, id)
		}
	}

	if len(payCenterIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid pay_center_ids provided"})
		return
	}

	// Check cache first for all IDs
	results := make(map[string]interface{})
	uncachedIDs := []int{}

	mapCacheMutex.RLock()
	for _, id := range payCenterIDs {
		idStr := strconv.Itoa(id)
		if cached, exists := mapCache[idStr]; exists {
			if time.Since(cached.timestamp) < cacheTimeout {
				results[idStr] = cached.data
			} else {
				// Remove expired cache
				delete(mapCache, idStr)
				uncachedIDs = append(uncachedIDs, id)
			}
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}
	mapCacheMutex.RUnlock()

	// Fetch uncached data in batch
	if len(uncachedIDs) > 0 {
		// Build batch query with placeholders
		placeholders := make([]string, len(uncachedIDs))
		args := make([]interface{}, len(uncachedIDs))
		for i, id := range uncachedIDs {
			placeholders[i] = fmt.Sprintf(":%d", i+1)
			args[i] = id
		}

		batchQuery := fmt.Sprintf(`
			SELECT 
				pcp.PAY_CENTER_ID,
				COUNT(DISTINCT pcp.OWNER_REGNO) as OWNER_COUNT,
				COUNT(*) as ACTIVITY_OPERATORS,
				NVL(SUM(pcp.PROPERTY_SIZE), 0) as AREA,
				NVL((SELECT COUNT(DISTINCT pm.MRCH_REGNO) 
				     FROM GPS.PAY_MARKET pm 
				     WHERE pm.PAY_CENTER_ID = pcp.PAY_CENTER_ID), 0) as TENANTS,
				NVL((SELECT SUM(
					NVL(lv.AREA_M2, 0)
				) FROM GPS.V_E_TUB_LAND_VIEW lv 
				WHERE lv.PAY_CENTER_ID = pcp.PAY_CENTER_ID), 0) as LAND_AREA
			FROM GPS.PAY_CENTER_PROPERTY pcp 
			WHERE pcp.PAY_CENTER_ID IN (%s)
			GROUP BY pcp.PAY_CENTER_ID`, strings.Join(placeholders, ","))

		rows, err := database.DB.Query(batchQuery, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
			return
		}
		defer rows.Close()

		mapCacheMutex.Lock()
		for rows.Next() {
			var payCenterID, ownerCount, activityOperators, tenants int
			var area, landArea float64

			if err := rows.Scan(&payCenterID, &ownerCount, &activityOperators, &area, &tenants, &landArea); err != nil {
				continue // Skip errors for individual rows
			}

			// Calculate unused area
			unusedArea := landArea - area

			mapData := map[string]interface{}{
				"owner_count":        ownerCount,
				"activity_operators": activityOperators,
				"area":               area,       // Ашиглагдаж байгаа талбай (мкв)
				"land_area":          landArea,   // Газрын талбай (мкв)
				"unused_area":        unusedArea, // Ашиглагдаагүй талбай (мкв)
				"tenants":            tenants,
			}

			idStr := strconv.Itoa(payCenterID)
			resultData := map[string]interface{}{
				"success": true,
				"data":    mapData,
			}
			results[idStr] = resultData

			// Cache the result
			mapCache[idStr] = &mapDataCache{
				data:      resultData,
				timestamp: time.Now(),
			}
		}
		mapCacheMutex.Unlock()

		// Add default data for IDs that weren't found in database
		for _, id := range uncachedIDs {
			idStr := strconv.Itoa(id)
			if _, exists := results[idStr]; !exists {
				defaultData := map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"owner_count":        0,
						"activity_operators": 0,
						"area":               0,
						"land_area":          0,
						"unused_area":        0,
						"tenants":            0,
					},
				}
				results[idStr] = defaultData

				// Cache the default result
				mapCacheMutex.Lock()
				mapCache[idStr] = &mapDataCache{
					data:      defaultData,
					timestamp: time.Now(),
				}
				mapCacheMutex.Unlock()
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// GetPayCenterPropertiesHandler returns all PAY_CENTER_PROPERTY data
func GetPayCenterPropertiesHandler(c *gin.Context) {
	query := `SELECT ID, PAY_CENTER_ID, PROPERTY_NUMBER, STATUS, CREATED_BY, CREATED_DATE, UPDATED_BY, UPDATED_DATE, PROPERTY_TYPE, OWNER_REGNO, PROPERTY_SIZE, RENT_AMOUNT FROM GPS.PAY_CENTER_PROPERTY`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, payCenterId, status, createdBy, updatedBy, propertyType *int
		var propertyNumber, createdDate, updatedDate, ownerRegno *string
		var propertySize, rentAmount *float64

		err := rows.Scan(&id, &payCenterId, &propertyNumber, &status, &createdBy, &createdDate, &updatedBy, &updatedDate, &propertyType, &ownerRegno, &propertySize, &rentAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"id":              id,
			"pay_center_id":   payCenterId,
			"property_number": propertyNumber,
			"status":          status,
			"created_by":      createdBy,
			"created_date":    createdDate,
			"updated_by":      updatedBy,
			"updated_date":    updatedDate,
			"property_type":   propertyType,
			"owner_regno":     ownerRegno,
			"property_size":   propertySize,
			"rent_amount":     rentAmount,
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
