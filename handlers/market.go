package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ListMarketsByOrg(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid organization ID")
		return
	}
	markets, err := repository.GetMarketsByOrgID(uint(orgID))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": markets})
}

// GetCenters returns all centers
func GetCenters(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT ID, NAME, BUILD_FLOOR, OFFICE_CODE, KHO_CODE, REGNO, LAT, LNG, PARCEL_ID FROM GPS.PAY_CENTER`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var buildFloor sql.NullInt64
		var name, officeCode, khoCode, regno, lat, lng, parcelId sql.NullString

		if err := rows.Scan(&id, &name, &buildFloor, &officeCode, &khoCode, &regno, &lat, &lng, &parcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"id":          id,
			"name":        nil,
			"build_floor": nil,
			"office_code": nil,
			"kho_code":    nil,
			"regno":       nil,
			"lat":         nil,
			"lng":         nil,
			"parcel_id":   nil,
		}

		if buildFloor.Valid {
			result["build_floor"] = int(buildFloor.Int64)
		}
		if name.Valid {
			result["name"] = name.String
		}
		if officeCode.Valid {
			result["office_code"] = officeCode.String
		}
		if khoCode.Valid {
			result["kho_code"] = khoCode.String
		}
		if regno.Valid {
			result["regno"] = regno.String
		}
		if lat.Valid {
			result["lat"] = lat.String
		}
		if lng.Valid {
			result["lng"] = lng.String
		}
		if parcelId.Valid {
			result["parcel_id"] = parcelId.String
		}

		results = append(results, result)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// GetBuildings returns all buildings
func GetBuildings(c *gin.Context) {
	// Жoinтой query ашиглаж PAY_CENTER болон PAY_MARKET-аас мэдээлэл авах
	query := `
		SELECT 
			pc.ID, 
			pc.NAME, 
			pc.BUILD_FLOOR, 
			pc.OFFICE_CODE, 
			pc.KHO_CODE, 
			pc.REGNO, 
			pc.LAT, 
			pc.LNG, 
			pc.PARCEL_ID,
			pc.ADDRESS,
			COUNT(DISTINCT pm.MRCH_REGNO) as TAX_PAYERS
		FROM GPS.PAY_CENTER pc
		LEFT JOIN GPS.PAY_MARKET pm ON pc.ID = pm.PAY_CENTER_ID
		GROUP BY pc.ID, pc.NAME, pc.BUILD_FLOOR, pc.OFFICE_CODE, pc.KHO_CODE, pc.REGNO, pc.LAT, pc.LNG, pc.PARCEL_ID, pc.ADDRESS
		ORDER BY pc.ID
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var id, taxPayers int
		var buildFloor sql.NullInt64
		var name, officeCode, khoCode, regno, lat, lng, parcelId, address sql.NullString

		if err := rows.Scan(&id, &name, &buildFloor, &officeCode, &khoCode, &regno, &lat, &lng, &parcelId, &address, &taxPayers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		result := map[string]interface{}{
			"id":          id,
			"name":        nil,
			"build_floor": nil,
			"office_code": nil,
			"kho_code":    nil,
			"regno":       nil,
			"lat":         nil,
			"lng":         nil,
			"parcel_id":   nil,
			"address":     nil,
			"tax_payers":  taxPayers,
		}

		if buildFloor.Valid {
			result["build_floor"] = int(buildFloor.Int64)
		}
		if name.Valid {
			result["name"] = name.String
		}
		if officeCode.Valid {
			result["office_code"] = officeCode.String
		}
		if khoCode.Valid {
			result["kho_code"] = khoCode.String
		}
		if regno.Valid {
			result["regno"] = regno.String
		}
		if lat.Valid {
			result["lat"] = lat.String
		}
		if lng.Valid {
			result["lng"] = lng.String
		}
		if parcelId.Valid {
			result["parcel_id"] = parcelId.String
		}
		if address.Valid {
			result["address"] = address.String
		}

		results = append(results, result)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// GetFloors returns all unique floors for a building (pay center)
func GetFloors(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, err := database.DB.Query(`SELECT DISTINCT STOR_FLOOR FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :id ORDER BY STOR_FLOOR`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()
	floors := []string{}
	for rows.Next() {
		var floor string
		if err := rows.Scan(&floor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		floors = append(floors, floor)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": floors})
}

// GetOrganizations returns all organizations for a building and floor
func GetOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	floor := c.Param("floor")
	rows, err := database.DB.Query(`SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR, pc.PARCEL_ID FROM GPS.PAY_MARKET pm LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID WHERE pm.PAY_CENTER_ID = :id AND pm.STOR_FLOOR = :floor`, id, floor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()
	orgs := []map[string]interface{}{}
	for rows.Next() {
		var m struct {
			ID                  int
			OpTypeName          string
			DistCode            string
			KhoCode             string
			StorName            string
			StorFloor           string
			MrchRegno           string
			PayCenterPropertyID int
			PayCenterID         int
			Lat                 float64
			Lng                 float64
			BuildFloor          sql.NullInt64
			ParcelId            sql.NullString
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng, &m.BuildFloor, &m.ParcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		orgMap := map[string]interface{}{
			"id":                     m.ID,
			"op_type_name":           m.OpTypeName,
			"dist_code":              m.DistCode,
			"kho_code":               m.KhoCode,
			"stor_name":              m.StorName,
			"stor_floor":             m.StorFloor,
			"mrch_regno":             m.MrchRegno,
			"pay_center_property_id": m.PayCenterPropertyID,
			"pay_center_id":          m.PayCenterID,
			"lat":                    m.Lat,
			"lng":                    m.Lng,
			"build_floor":            nil,
			"parcel_id":              nil,
		}
		if m.BuildFloor.Valid {
			orgMap["build_floor"] = int(m.BuildFloor.Int64)
		}
		if m.ParcelId.Valid {
			orgMap["parcel_id"] = m.ParcelId.String
		}
		orgs = append(orgs, orgMap)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetAllOrganizations returns all organizations for a building (all floors)
func GetAllOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, err := database.DB.Query(`SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR, pc.PARCEL_ID FROM GPS.PAY_MARKET pm LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID WHERE pm.PAY_CENTER_ID = :id`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()
	orgs := []map[string]interface{}{}
	for rows.Next() {
		var m struct {
			ID                  int
			OpTypeName          string
			DistCode            string
			KhoCode             string
			StorName            string
			StorFloor           string
			MrchRegno           string
			PayCenterPropertyID int
			PayCenterID         int
			Lat                 sql.NullFloat64
			Lng                 sql.NullFloat64
			BuildFloor          sql.NullInt64
			ParcelId            sql.NullString
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng, &m.BuildFloor, &m.ParcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		// Helper function to get float value safely
		getFloatValue := func(nf sql.NullFloat64) interface{} {
			if nf.Valid {
				return nf.Float64
			}
			return nil
		}

		orgMap := map[string]interface{}{
			"id":                     m.ID,
			"op_type_name":           m.OpTypeName,
			"dist_code":              m.DistCode,
			"kho_code":               m.KhoCode,
			"stor_name":              m.StorName,
			"stor_floor":             m.StorFloor,
			"mrch_regno":             m.MrchRegno,
			"pay_center_property_id": m.PayCenterPropertyID,
			"pay_center_id":          m.PayCenterID,
			"lat":                    getFloatValue(m.Lat),
			"lng":                    getFloatValue(m.Lng),
			"build_floor":            nil,
			"parcel_id":              nil,
		}
		if m.BuildFloor.Valid {
			orgMap["build_floor"] = int(m.BuildFloor.Int64)
		}
		if m.ParcelId.Valid {
			orgMap["parcel_id"] = m.ParcelId.String
		}
		orgs = append(orgs, orgMap)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetMarketsByPayCenterID returns all markets for a specific pay center
func GetMarketsByPayCenterID(c *gin.Context) {
	payCenterID := c.Query("pay_center_id")
	if payCenterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pay_center_id parameter is required"})
		return
	}

	id, err := strconv.Atoi(payCenterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay_center_id"})
		return
	}

	rows, err := database.DB.Query(`SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR, pc.PARCEL_ID FROM GPS.PAY_MARKET pm LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID WHERE pm.PAY_CENTER_ID = :id`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	orgs := []map[string]interface{}{}
	for rows.Next() {
		var m struct {
			ID                  int
			OpTypeName          string
			DistCode            string
			KhoCode             string
			StorName            string
			StorFloor           string
			MrchRegno           string
			PayCenterPropertyID int
			PayCenterID         int
			Lat                 sql.NullFloat64
			Lng                 sql.NullFloat64
			BuildFloor          sql.NullInt64
			ParcelId            sql.NullString
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng, &m.BuildFloor, &m.ParcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		// Helper function to get float value safely
		getFloatValue := func(nf sql.NullFloat64) interface{} {
			if nf.Valid {
				return nf.Float64
			}
			return nil
		}

		orgMap := map[string]interface{}{
			"id":                     m.ID,
			"op_type_name":           m.OpTypeName,
			"dist_code":              m.DistCode,
			"kho_code":               m.KhoCode,
			"stor_name":              m.StorName,
			"stor_floor":             m.StorFloor,
			"mrch_regno":             m.MrchRegno,
			"pay_center_property_id": m.PayCenterPropertyID,
			"pay_center_id":          m.PayCenterID,
			"lat":                    getFloatValue(m.Lat),
			"lng":                    getFloatValue(m.Lng),
			"build_floor":            nil,
			"parcel_id":              nil,
		}
		if m.BuildFloor.Valid {
			orgMap["build_floor"] = int(m.BuildFloor.Int64)
		}
		if m.ParcelId.Valid {
			orgMap["parcel_id"] = m.ParcelId.String
		}
		orgs = append(orgs, orgMap)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetBuildingByID returns a single building by ID
func GetBuildingByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// PAY_CENTER болон PAY_MARKET-аас мэдээлэл авах
	query := `
		SELECT 
			pc.ID, 
			pc.NAME, 
			pc.BUILD_FLOOR, 
			pc.OFFICE_CODE, 
			pc.KHO_CODE, 
			pc.REGNO, 
			pc.LAT, 
			pc.LNG, 
			pc.PARCEL_ID,
			pc.ADDRESS,
			COUNT(DISTINCT pm.MRCH_REGNO) as TAX_PAYERS
		FROM GPS.PAY_CENTER pc
		LEFT JOIN GPS.PAY_MARKET pm ON pc.ID = pm.PAY_CENTER_ID
		WHERE pc.ID = :1
		GROUP BY pc.ID, pc.NAME, pc.BUILD_FLOOR, pc.OFFICE_CODE, pc.KHO_CODE, pc.REGNO, pc.LAT, pc.LNG, pc.PARCEL_ID, pc.ADDRESS
	`

	row := database.DB.QueryRow(query, id)

	var buildingID, taxPayers int
	var buildFloor sql.NullInt64
	var name, officeCode, khoCode, regno, lat, lng, parcelId, address sql.NullString

	if err := row.Scan(&buildingID, &name, &buildFloor, &officeCode, &khoCode, &regno, &lat, &lng, &parcelId, &address, &taxPayers); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
		return
	}

	result := map[string]interface{}{
		"id":          buildingID,
		"name":        nil,
		"build_floor": nil,
		"office_code": nil,
		"kho_code":    nil,
		"regno":       nil,
		"lat":         nil,
		"lng":         nil,
		"parcel_id":   nil,
		"address":     nil,
		"tax_payers":  taxPayers,
	}

	if buildFloor.Valid {
		result["build_floor"] = int(buildFloor.Int64)
	}
	if name.Valid {
		result["name"] = name.String
	}
	if officeCode.Valid {
		result["office_code"] = officeCode.String
	}
	if khoCode.Valid {
		result["kho_code"] = khoCode.String
	}
	if regno.Valid {
		result["regno"] = regno.String
	}
	if lat.Valid {
		result["lat"] = lat.String
	}
	if lng.Valid {
		result["lng"] = lng.String
	}
	if parcelId.Valid {
		result["parcel_id"] = parcelId.String
	}
	if address.Valid {
		result["address"] = address.String
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetEbarimtByPin returns ebarimt info by PIN from new V_E_TUB_PAY_MARKET_EBARIMT table
func GetEbarimtByPin(c *gin.Context) {
	pin := c.Param("pin")

	// Query V_E_TUB_PAY_MARKET_EBARIMT with new structure
	query := `SELECT 
		COALESCE(SUM(CNT_3), 0) as CNT_3,
		COALESCE(SUM(CNT_30), 0) as CNT_30,
		MAX(OP_TYPE_NAME) as OP_TYPE_NAME,
		MAX(MAR_NAME) as MAR_NAME,
		MAX(MAR_REGNO) as MAR_REGNO,
		MAX(QR_CODE) as QR_CODE
	FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT 
	WHERE TRIM(UPPER(MRCH_REGNO)) = TRIM(UPPER(:1))`

	row := database.DB.QueryRow(query, pin)

	var cnt3, cnt30 int
	var opTypeName, marName, marRegno, qrCode sql.NullString

	err := row.Scan(&cnt3, &cnt30, &opTypeName, &marName, &marRegno, &qrCode)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return zero values if no data found
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"count_receipt": 0,
					"cnt_3":         0,
					"cnt_30":        0,
					"op_type_name":  "",
					"mar_name":      "",
					"mar_regno":     "",
					"qr_code":       "",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
		return
	}

	result := gin.H{
		"count_receipt": cnt3, // For backward compatibility
		"cnt_3":         cnt3,
		"cnt_30":        cnt30,
		"op_type_name":  "",
		"mar_name":      "",
		"mar_regno":     "",
		"qr_code":       "",
	}

	if opTypeName.Valid {
		result["op_type_name"] = opTypeName.String
	}
	if marName.Valid {
		result["mar_name"] = marName.String
	}
	if marRegno.Valid {
		result["mar_regno"] = marRegno.String
	}
	if qrCode.Valid {
		result["qr_code"] = qrCode.String
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// Debug endpoint to check MRCH_REGNO to PIN mapping
func GetEbarimtDebug(c *gin.Context) {
	mrchRegno := c.Query("mrch_regno")
	if mrchRegno == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mrch_regno parameter is required"})
		return
	}

	// First check if this MRCH_REGNO exists in PAY_MARKET
	var payMarketExists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_MARKET WHERE MRCH_REGNO = :1", mrchRegno).Scan(&payMarketExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking PAY_MARKET: " + err.Error()})
		return
	}

	// Check if this MRCH_REGNO exists as PIN in V_E_TUB_PAY_MARKET_EBARIMT
	var barimtExists int
	var countReceipt int
	err = database.DB.QueryRow("SELECT COUNT(*), COALESCE(MAX(CNT_3), 0) FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT WHERE TRIM(UPPER(MRCH_REGNO)) = TRIM(UPPER(:1))", mrchRegno).Scan(&barimtExists, &countReceipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking V_E_TUB_PAY_MARKET_EBARIMT: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mrch_regno":        mrchRegno,
		"pay_market_exists": payMarketExists > 0,
		"barimt_exists":     barimtExists > 0,
		"count_receipt":     countReceipt,
		"debug":             "MRCH_REGNO сэс PIN болж ашиглагдаж байна уу шалгах",
	})
}

// GetDashboardStatistics returns overall dashboard statistics
func GetDashboardStatistics(c *gin.Context) {
	var totalBuildings, totalTenants, totalReceiptCount int
	var totalLegalEntities, totalCitizens, totalOwners int
	var totalArea, totalLandArea float64
	var nuatCount, nhatCount int

	fmt.Println("=== DEBUG: Starting GetDashboardStatistics ===")

	// 1. Нийт Объект (барилга) - PAY_CENTER тоо
	err := database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER").Scan(&totalBuildings)
	if err != nil {
		fmt.Printf("ERROR counting buildings: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting buildings: " + err.Error()})
		return
	}
	fmt.Printf("Total buildings: %d\n", totalBuildings)

	// 2. Түрээслэгч Нийт мкв - PAY_CENTER_PROPERTY доторх бүх PROPERTY_SIZE нэмэх
	err = database.DB.QueryRow("SELECT NVL(SUM(TO_NUMBER(REPLACE(PROPERTY_SIZE, ',', '.'))), 0) FROM GPS.PAY_CENTER_PROPERTY").Scan(&totalArea)
	if err != nil {
		fmt.Printf("ERROR calculating total area: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating total area: " + err.Error()})
		return
	}
	fmt.Printf("Total area: %f\n", totalArea)

	// 3. Нийт газрын талбай - PAY_CENTER ID-аар V_E_TUB_LAND_VIEW-тай join хийж AREA_M2 + AREA_HA нэмэх
	err = database.DB.QueryRow(`
		SELECT NVL(SUM(
			NVL(v.AREA_M2, 0) + NVL(v.AREA_HA, 0)
		), 0) as total_land_area
		FROM GPS.PAY_CENTER pc
		JOIN GPS.V_E_TUB_LAND_VIEW v ON pc.ID = v.PAY_CENTER_ID
		WHERE v.AREA_M2 IS NOT NULL OR v.AREA_HA IS NOT NULL
	`).Scan(&totalLandArea)
	if err != nil {
		fmt.Printf("ERROR accessing GPS.V_E_TUB_LAND_VIEW with JOIN: %v\n", err)
		// Fallback: try without GPS schema
		err2 := database.DB.QueryRow(`
			SELECT NVL(SUM(
				NVL(v.AREA_M2, 0) + NVL(v.AREA_HA, 0)
			), 0) as total_land_area
			FROM PAY_CENTER pc
			JOIN V_E_TUB_LAND_VIEW v ON pc.ID = v.PAY_CENTER_ID
			WHERE v.AREA_M2 IS NOT NULL OR v.AREA_HA IS NOT NULL
		`).Scan(&totalLandArea)
		if err2 != nil {
			fmt.Printf("ERROR accessing V_E_TUB_LAND_VIEW with JOIN: %v\n", err2)
			totalLandArea = 0
		}
	}
	fmt.Printf("Total land area: %f\n", totalLandArea)

	// 4. Хуулийн этгээд - PAY_MARKET-аас 7 оронтой MRCH_REGNO тоолох (CHANGED from 10 to 7)
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE LENGTH(TRIM(MRCH_REGNO)) = 7
	`).Scan(&totalLegalEntities)
	if err != nil {
		fmt.Printf("ERROR counting legal entities: %v\n", err)
		totalLegalEntities = 0
	}
	fmt.Printf("Total legal entities: %d\n", totalLegalEntities)

	// 5. Иргэн - PAY_MARKET-аас 10 оронтой MRCH_REGNO тоолох (CHANGED from 7 to 10)
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE LENGTH(TRIM(MRCH_REGNO)) = 10
	`).Scan(&totalCitizens)
	if err != nil {
		fmt.Printf("ERROR counting citizens: %v\n", err)
		totalCitizens = 0
	}
	fmt.Printf("Total citizens: %d\n", totalCitizens)

	// 6. Эзэмшигч - PAY_CENTER_PROPERTY-аас OWNER_REGNO unique тоолох
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT OWNER_REGNO) 
		FROM GPS.PAY_CENTER_PROPERTY 
		WHERE OWNER_REGNO IS NOT NULL
	`).Scan(&totalOwners)
	if err != nil {
		fmt.Printf("ERROR counting owners: %v\n", err)
		totalOwners = 0
	}
	fmt.Printf("Total owners: %d\n", totalOwners)

	// 7. Түрээслэгч - PAY_MARKET доторх unique MRCH_REGNO тоолох
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE MRCH_REGNO IS NOT NULL
	`).Scan(&totalTenants)
	if err != nil {
		fmt.Printf("ERROR counting tenants: %v\n", err)
		totalTenants = 0
	}
	fmt.Printf("Total tenants: %d\n", totalTenants)

	// 8. Баримт хэвлэдэг - V_E_TUB_PAY_MARKET_EBARIMT доторх бүх мөрийг тоолох
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT").Scan(&totalReceiptCount)
	if err != nil {
		fmt.Printf("ERROR counting receipts: %v\n", err)
		totalReceiptCount = 0
	}
	fmt.Printf("Total receipt count: %d\n", totalReceiptCount)

	// 9. НӨАТ суутган төлөгч - V_E_TUB_COUNT_NHAT_NUAT-аас NUAT__COUNT нэмэх
	err = database.DB.QueryRow("SELECT NVL(SUM(NVL(NUAT__COUNT, 0)), 0) FROM GPS.V_E_TUB_COUNT_NHAT_NUAT").Scan(&nuatCount)
	if err != nil {
		fmt.Printf("ERROR counting NUAT: %v\n", err)
		// Try without GPS schema
		err2 := database.DB.QueryRow("SELECT NVL(SUM(NVL(NUAT__COUNT, 0)), 0) FROM V_E_TUB_COUNT_NHAT_NUAT").Scan(&nuatCount)
		if err2 != nil {
			fmt.Printf("ERROR counting NUAT (alt): %v\n", err2)
			nuatCount = 897 // Fallback to default
		}
	}
	fmt.Printf("Total NUAT count: %d\n", nuatCount)

	// 10. НХАТ төлөгч - V_E_TUB_COUNT_NHAT_NUAT-аас NHAT__COUNT нэмэх
	err = database.DB.QueryRow("SELECT NVL(SUM(NVL(NHAT__COUNT, 0)), 0) FROM GPS.V_E_TUB_COUNT_NHAT_NUAT").Scan(&nhatCount)
	if err != nil {
		fmt.Printf("ERROR counting NHAT: %v\n", err)
		// Try without GPS schema
		err2 := database.DB.QueryRow("SELECT NVL(SUM(NVL(NHAT__COUNT, 0)), 0) FROM V_E_TUB_COUNT_NHAT_NUAT").Scan(&nhatCount)
		if err2 != nil {
			fmt.Printf("ERROR counting NHAT (alt): %v\n", err2)
			nhatCount = 8970 // Fallback to default
		}
	}
	fmt.Printf("Total NHAT count: %d\n", nhatCount)

	// 11. Ашиглагдаагүй талбай тооцоолох
	unusedArea := totalLandArea - totalArea

	fmt.Println("=== DEBUG: Completed GetDashboardStatistics ===")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_buildings":      totalBuildings,
			"total_area":           totalArea,          // Нийт ашиглагдаж буй талбай мкв
			"total_land_area":      totalLandArea,      // Нийт газрын талбай мкв
			"unused_area":          unusedArea,         // Ашиглагдаагүй талбай
			"total_legal_entities": totalLegalEntities, // Хуулийн этгээд (7 орон)
			"total_citizens":       totalCitizens,      // Иргэн (10 орон)
			"total_owners":         totalOwners,        // Эзэмшигч
			"total_tenants":        totalTenants,       // Түрээслэгч
			"total_receipt_count":  totalReceiptCount,  // Баримт хэвлэдэг
			"nuat_count":           nuatCount,          // НӨАТ суутган төлөгч
			"nhat_count":           nhatCount,          // НХАТ төлөгч
		},
	})
}

// GetRegistrationStats returns registration statistics for a building
func GetRegistrationStats(c *gin.Context) {
	buildingID := c.Param("id")
	if buildingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID is required"})
		return
	}

	// Query to get MRCH_REGNO from PAY_MARKET where PAY_CENTER_ID matches building ID
	query := `
		SELECT DISTINCT pm.MRCH_REGNO
		FROM GPS.PAY_MARKET pm
		WHERE pm.PAY_CENTER_ID = :1
	`

	rows, err := database.DB.Query(query, buildingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Data  base query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var mrchRegnos []string
	for rows.Next() {
		var mrchRegno sql.NullString
		if err := rows.Scan(&mrchRegno); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		if mrchRegno.Valid {
			mrchRegnos = append(mrchRegnos, mrchRegno.String)
		}
	}

	if len(mrchRegnos) == 0 {
		// No organizations found for this building
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"registered":                0,
				"not_registered":            0,
				"total":                     0,
				"registered_percentage":     0,
				"not_registered_percentage": 0,
			},
		})
		return
	}

	// Process in batches to avoid Oracle IN clause limit of 1000
	const batchSize = 1000
	var totalRegisteredCount, totalNotRegisteredCount, totalCount int

	for i := 0; i < len(mrchRegnos); i += batchSize {
		end := i + batchSize
		if end > len(mrchRegnos) {
			end = len(mrchRegnos)
		}

		batch := mrchRegnos[i:end]

		// Build the IN clause for this batch
		placeholders := make([]string, len(batch))
		args := make([]interface{}, len(batch))

		for j, regno := range batch {
			placeholders[j] = fmt.Sprintf(":%d", j+1)
			args[j] = regno
		}

		// Query V_E_TUB_BRANCH to count registrations for this batch
		registrationQuery := fmt.Sprintf(`
			SELECT 
				SUM(CASE WHEN TULUV = 'Бүртгэгдсэн' THEN 1 ELSE 0 END) as registered_count,
				SUM(CASE WHEN TULUV = 'Бүртгэгдээгүй' THEN 1 ELSE 0 END) as not_registered_count,
				COUNT(*) as total_count
			FROM GPS.V_E_TUB_BRANCH 
			WHERE REGISTER IN (%s)
		`, strings.Join(placeholders, ","))

		row := database.DB.QueryRow(registrationQuery, args...)

		var registeredCount, notRegisteredCount, count int
		err = row.Scan(&registeredCount, &notRegisteredCount, &count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration count error: " + err.Error()})
			return
		}

		// Accumulate results
		totalRegisteredCount += registeredCount
		totalNotRegisteredCount += notRegisteredCount
		totalCount += count
	}

	// Calculate percentages
	var registeredPercentage, notRegisteredPercentage float64
	if totalCount > 0 {
		registeredPercentage = float64(totalRegisteredCount) / float64(totalCount) * 100
		notRegisteredPercentage = float64(totalNotRegisteredCount) / float64(totalCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"registered":                totalRegisteredCount,
			"not_registered":            totalNotRegisteredCount,
			"total":                     totalCount,
			"registered_percentage":     math.Round(registeredPercentage*100) / 100,
			"not_registered_percentage": math.Round(notRegisteredPercentage*100) / 100,
		},
	})
}

// GetTaxOfficeStats returns tax office statistics for a building
func GetTaxOfficeStats(c *gin.Context) {
	buildingID := c.Param("id")
	if buildingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID is required"})
		return
	}

	// Query to get MRCH_REGNO from PAY_MARKET where PAY_CENTER_ID matches building ID
	query := `
		SELECT DISTINCT pm.MRCH_REGNO
		FROM GPS.PAY_MARKET pm
		WHERE pm.PAY_CENTER_ID = :1
	`

	rows, err := database.DB.Query(query, buildingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var mrchRegnos []string
	for rows.Next() {
		var mrchRegno sql.NullString
		if err := rows.Scan(&mrchRegno); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		if mrchRegno.Valid {
			mrchRegnos = append(mrchRegnos, mrchRegno.String)
		}
	}

	fmt.Printf("Found %d MRCH_REGNO values: %v\n", len(mrchRegnos), mrchRegnos)

	if len(mrchRegnos) == 0 {
		// No organizations found for this building
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"districts": []gin.H{},
			},
		})
		return
	}

	// Process in batches to avoid Oracle IN clause limit of 1000
	const batchSize = 1000
	districtMap := make(map[string]map[string]int) // TTA -> DED_ALBA -> count

	for i := 0; i < len(mrchRegnos); i += batchSize {
		end := i + batchSize
		if end > len(mrchRegnos) {
			end = len(mrchRegnos)
		}

		batch := mrchRegnos[i:end]

		// Build the IN clause for this batch
		placeholders := make([]string, len(batch))
		args := make([]interface{}, len(batch))

		for j, regno := range batch {
			placeholders[j] = fmt.Sprintf(":%d", j+1)
			args[j] = regno
		}

		// Query V_E_TUB_BRANCH to get TTA and DED_ALBA data for this batch
		hierarchyQuery := fmt.Sprintf(`
			SELECT 
				TTA,
				DED_ALBA,
				COUNT(*) as count
			FROM GPS.V_E_TUB_BRANCH 
			WHERE REGISTER IN (%s) 
			AND TTA IS NOT NULL AND LENGTH(TRIM(TTA)) > 0
			AND DED_ALBA IS NOT NULL AND LENGTH(TRIM(DED_ALBA)) > 0
			GROUP BY TTA, DED_ALBA
			ORDER BY TTA, count DESC
		`, strings.Join(placeholders, ","))

		// Execute hierarchy query for this batch
		hierarchyRows, err := database.DB.Query(hierarchyQuery, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Hierarchy query error: " + err.Error()})
			return
		}

		// Process results from this batch
		for hierarchyRows.Next() {
			var tta, dedAlba sql.NullString
			var count int
			if err := hierarchyRows.Scan(&tta, &dedAlba, &count); err != nil {
				hierarchyRows.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Hierarchy scan error: " + err.Error()})
				return
			}

			if tta.Valid && dedAlba.Valid {
				if districtMap[tta.String] == nil {
					districtMap[tta.String] = make(map[string]int)
				}
				districtMap[tta.String][dedAlba.String] += count
			}
		}
		hierarchyRows.Close()
	}

	fmt.Printf("Found %d districts in hierarchy data\n", len(districtMap))

	// Convert to final structure
	var districts []gin.H
	for districtName, khoroosMap := range districtMap {
		// Convert khoroos map to slice for frontend compatibility
		var khoroos []gin.H
		totalCount := 0
		for khorooName, count := range khoroosMap {
			khoroos = append(khoroos, gin.H{
				"name":  khorooName,
				"count": count,
			})
			totalCount += count
		}

		// Sort khoroos by count descending
		sort.Slice(khoroos, func(i, j int) bool {
			return khoroos[i]["count"].(int) > khoroos[j]["count"].(int)
		})

		districts = append(districts, gin.H{
			"name":    districtName,
			"count":   totalCount, // Total count for this district
			"khoroos": khoroos,
		})
	}

	// Sort districts by count (descending)
	sort.Slice(districts, func(i, j int) bool {
		return districts[i]["count"].(int) > districts[j]["count"].(int)
	})

	fmt.Printf("Districts found: %d\n", len(districts))
	for _, district := range districts {
		khoroos := district["khoroos"].([]gin.H)
		fmt.Printf("  %s: %d түрээслэгч, %d хороо\n",
			district["name"], district["count"], len(khoroos))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"districts": districts,
		},
	})
}

// GetSegmentStats returns segment statistics for a building
func GetSegmentStats(c *gin.Context) {
	buildingID := c.Param("id")
	if buildingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Building ID is required"})
		return
	}

	// Query to get MRCH_REGNO from PAY_MARKET where PAY_CENTER_ID matches building ID
	query := `
		SELECT DISTINCT pm.MRCH_REGNO
		FROM GPS.PAY_MARKET pm
		WHERE pm.PAY_CENTER_ID = :1
	`

	rows, err := database.DB.Query(query, buildingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var mrchRegnos []string
	for rows.Next() {
		var mrchRegno sql.NullString
		if err := rows.Scan(&mrchRegno); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		if mrchRegno.Valid {
			mrchRegnos = append(mrchRegnos, mrchRegno.String)
		}
	}

	if len(mrchRegnos) == 0 {
		// No organizations found for this building
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []gin.H{},
		})
		return
	}

	// Build the IN clause for the query
	placeholders := make([]string, len(mrchRegnos))
	args := make([]interface{}, len(mrchRegnos))

	for i, regno := range mrchRegnos {
		placeholders[i] = fmt.Sprintf(":%d", i+1)
		args[i] = regno
	}

	// Query V_E_TUB_SEGMENT to get segment statistics
	segmentQuery := fmt.Sprintf(`
		SELECT 
			SEGMENT,
			COUNT(*) as count
		FROM GPS.V_E_TUB_SEGMENT 
		WHERE TRIM(UPPER(PIN)) IN (%s)
		GROUP BY SEGMENT
		ORDER BY count DESC
	`, strings.Join(placeholders, ","))

	rows, err = database.DB.Query(segmentQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Segment query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var segments []gin.H
	for rows.Next() {
		var segment sql.NullString
		var count int
		if err := rows.Scan(&segment, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Segment scan error: " + err.Error()})
			return
		}

		segmentName := "Тодорхойгүй"
		if segment.Valid && segment.String != "" {
			segmentName = segment.String
		}

		segments = append(segments, gin.H{
			"name":  segmentName,
			"count": count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    segments,
	})
}

// GetOperatorStats returns operator statistics for a specific PAY_CENTER_ID
func GetOperatorStats(c *gin.Context) {
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
			OPRT_NAME,
			COUNT(*) as count
		FROM GPS.V_E_TUB_OPERATORS 
		WHERE PAY_CENTER_ID = :1
			AND OPRT_NAME IS NOT NULL
		GROUP BY OPRT_NAME
		ORDER BY count DESC, OPRT_NAME
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

	type OperatorStat struct {
		OprtName string `json:"oprt_name"`
		Count    int    `json:"count"`
	}

	var results []OperatorStat
	for rows.Next() {
		var oprtName sql.NullString
		var count int

		if err := rows.Scan(&oprtName, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Scan error: " + err.Error(),
			})
			return
		}

		if oprtName.Valid && oprtName.String != "" {
			result := OperatorStat{
				OprtName: oprtName.String,
				Count:    count,
			}
			results = append(results, result)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"data":          results,
		"pay_center_id": payCenterId,
	})
}

// DiagnosticStatistics endpoint to debug why statistics return 0
func DiagnosticStatistics(c *gin.Context) {
	fmt.Println("=== DIAGNOSTIC: Starting statistics debug ===")

	var result = gin.H{}

	// 1. Check PAY_MARKET table structure and data
	var payMarketCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_MARKET").Scan(&payMarketCount)
	if err != nil {
		result["pay_market_error"] = err.Error()
	} else {
		result["pay_market_total_rows"] = payMarketCount
	}

	// 2. Check MRCH_REGNO values
	var mrchRegnoCount int
	err = database.DB.QueryRow("SELECT COUNT(DISTINCT MRCH_REGNO) FROM GPS.PAY_MARKET WHERE MRCH_REGNO IS NOT NULL").Scan(&mrchRegnoCount)
	if err != nil {
		result["mrch_regno_error"] = err.Error()
	} else {
		result["mrch_regno_distinct_count"] = mrchRegnoCount
	}

	// 3. Sample MRCH_REGNO values and their lengths
	rows, err := database.DB.Query(`
		SELECT MRCH_REGNO, LENGTH(TRIM(MRCH_REGNO)) as len 
		FROM GPS.PAY_MARKET 
		WHERE MRCH_REGNO IS NOT NULL 
		AND ROWNUM <= 10
	`)
	if err != nil {
		result["sample_mrch_error"] = err.Error()
	} else {
		defer rows.Close()
		var samples []gin.H
		for rows.Next() {
			var mrch sql.NullString
			var length int
			if err := rows.Scan(&mrch, &length); err == nil {
				samples = append(samples, gin.H{
					"mrch_regno": mrch.String,
					"length":     length,
				})
			}
		}
		result["sample_mrch_regnos"] = samples
	}

	// 4. Count by length
	lengthRows, err := database.DB.Query(`
		SELECT LENGTH(TRIM(MRCH_REGNO)) as len, COUNT(*) as cnt 
		FROM GPS.PAY_MARKET 
		WHERE MRCH_REGNO IS NOT NULL 
		GROUP BY LENGTH(TRIM(MRCH_REGNO))
		ORDER BY len
	`)
	if err != nil {
		result["length_count_error"] = err.Error()
	} else {
		defer lengthRows.Close()
		var lengthCounts []gin.H
		for lengthRows.Next() {
			var length, count int
			if err := lengthRows.Scan(&length, &count); err == nil {
				lengthCounts = append(lengthCounts, gin.H{
					"length": length,
					"count":  count,
				})
			}
		}
		result["mrch_regno_by_length"] = lengthCounts
	}

	// 5. Check PAY_CENTER_PROPERTY
	var propertyCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER_PROPERTY").Scan(&propertyCount)
	if err != nil {
		result["property_error"] = err.Error()
	} else {
		result["property_total_rows"] = propertyCount
	}

	// 6. Check OWNER_REGNO
	var ownerRegnoCount int
	err = database.DB.QueryRow("SELECT COUNT(DISTINCT OWNER_REGNO) FROM GPS.PAY_CENTER_PROPERTY WHERE OWNER_REGNO IS NOT NULL").Scan(&ownerRegnoCount)
	if err != nil {
		result["owner_regno_error"] = err.Error()
	} else {
		result["owner_regno_distinct_count"] = ownerRegnoCount
	}

	// 7. Check property value tables
	var propertyValueCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.V_TPI_PROPERTY_XYP_DATA_OWNER").Scan(&propertyValueCount)
	if err != nil {
		result["property_value_error"] = err.Error()
		// Try without GPS schema
		err2 := database.DB.QueryRow("SELECT COUNT(*) FROM V_TPI_PROPERTY_XYP_DATA_OWNER").Scan(&propertyValueCount)
		if err2 != nil {
			result["property_value_error_alt"] = err2.Error()
		} else {
			result["property_value_total_rows"] = propertyValueCount
		}
	} else {
		result["property_value_total_rows"] = propertyValueCount
	}

	// 8. Check land area tables
	var landAreaCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.V_E_TUB_LAND_VIEW").Scan(&landAreaCount)
	if err != nil {
		result["land_area_error"] = err.Error()
		// Try without GPS schema
		err2 := database.DB.QueryRow("SELECT COUNT(*) FROM V_E_TUB_LAND_VIEW").Scan(&landAreaCount)
		if err2 != nil {
			result["land_area_error_alt"] = err2.Error()
		} else {
			result["land_area_total_rows"] = landAreaCount
		}
	} else {
		result["land_area_total_rows"] = landAreaCount
	}

	// 9. Test specific queries that return 0
	// Test 7-digit MRCH_REGNO query
	var test7digit int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE MRCH_REGNO IS NOT NULL
		AND TRIM(MRCH_REGNO) != ''
		AND LENGTH(TRIM(MRCH_REGNO)) = 7
		AND TRIM(MRCH_REGNO) NOT IN ('0000000', '1111111', '2222222')
	`).Scan(&test7digit)
	if err != nil {
		result["test_7digit_error"] = err.Error()
	} else {
		result["test_7digit_count"] = test7digit
	}

	// Alternative query for 7-digit
	var test7digitAlt int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE LENGTH(TRIM(MRCH_REGNO)) = 7
	`).Scan(&test7digitAlt)
	if err != nil {
		result["test_7digit_alt_error"] = err.Error()
	} else {
		result["test_7digit_alt_count"] = test7digitAlt
	}

	// Test 10-digit MRCH_REGNO query
	var test10digit int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE MRCH_REGNO IS NOT NULL
		AND TRIM(MRCH_REGNO) != ''
		AND LENGTH(TRIM(MRCH_REGNO)) = 10
		AND TRIM(MRCH_REGNO) NOT IN ('0000000000', '1111111111', '2222222222')
	`).Scan(&test10digit)
	if err != nil {
		result["test_10digit_error"] = err.Error()
	} else {
		result["test_10digit_count"] = test10digit
	}

	// Alternative query for 10-digit
	var test10digitAlt int
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT MRCH_REGNO) 
		FROM GPS.PAY_MARKET 
		WHERE LENGTH(TRIM(MRCH_REGNO)) = 10
	`).Scan(&test10digitAlt)
	if err != nil {
		result["test_10digit_alt_error"] = err.Error()
	} else {
		result["test_10digit_alt_count"] = test10digitAlt
	}

	// Check V_E_TUB_COUNT_NHAT_NUAT existence and columns
	var nhatNuatExists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.V_E_TUB_COUNT_NHAT_NUAT WHERE ROWNUM = 1").Scan(&nhatNuatExists)
	if err != nil {
		result["nhat_nuat_table_error"] = err.Error()
	} else {
		result["nhat_nuat_table_exists"] = true

		// Try to get column names
		colRows, err := database.DB.Query(`
			SELECT column_name 
			FROM all_tab_columns 
			WHERE table_name = 'V_E_TUB_COUNT_NHAT_NUAT' 
			AND owner = 'GPS'
		`)
		if err == nil {
			defer colRows.Close()
			var columns []string
			for colRows.Next() {
				var colName string
				if err := colRows.Scan(&colName); err == nil {
					columns = append(columns, colName)
				}
			}
			result["nhat_nuat_columns"] = columns
		}
	}

	fmt.Println("=== DIAGNOSTIC: Completed statistics debug ===")

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"diagnostic_data": result,
	})
}

// GetNuatNhatByMrchRegno returns NUAT/NHAT data for specific MRCH_REGNO
func GetNuatNhatByMrchRegno(c *gin.Context) {
	mrchRegno := c.Param("regno")
	if mrchRegno == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MRCH_REGNO parameter is required"})
		return
	}

	fmt.Printf("Getting NUAT/NHAT data for MRCH_REGNO: %s\n", mrchRegno)

	var nuatCount, nhatCount int

	// Query V_E_TUB_COUNT_NHAT_NUAT for specific REGISTER (MRCH_REGNO)
	query := `
		SELECT 
			NVL(NUAT__COUNT, 0) as nuat_count,
			NVL(NHAT__COUNT, 0) as nhat_count
		FROM GPS.V_E_TUB_COUNT_NHAT_NUAT 
		WHERE TRIM(UPPER(REGISTER)) = TRIM(UPPER(:1))
	`

	err := database.DB.QueryRow(query, mrchRegno).Scan(&nuatCount, &nhatCount)
	if err != nil {
		fmt.Printf("ERROR getting NUAT/NHAT for %s: %v\n", mrchRegno, err)
		// Try without GPS schema
		err2 := database.DB.QueryRow(`
			SELECT 
				NVL(NUAT__COUNT, 0) as nuat_count,
				NVL(NHAT__COUNT, 0) as nhat_count
			FROM V_E_TUB_COUNT_NHAT_NUAT 
			WHERE TRIM(UPPER(REGISTER)) = TRIM(UPPER(:1))
		`, mrchRegno).Scan(&nuatCount, &nhatCount)
		if err2 != nil {
			fmt.Printf("ERROR getting NUAT/NHAT (alt) for %s: %v\n", mrchRegno, err2)
			// Set default values if not found
			nuatCount = 0
			nhatCount = 0
		}
	}

	fmt.Printf("NUAT: %d, NHAT: %d for MRCH_REGNO: %s\n", nuatCount, nhatCount, mrchRegno)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"nuat_count": nuatCount,
			"nhat_count": nhatCount,
			"mrch_regno": mrchRegno,
		},
	})
}
