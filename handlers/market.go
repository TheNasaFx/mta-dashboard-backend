package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"database/sql"
	"net/http"
	"strconv"

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
	rows, err := database.DB.Query(`SELECT DISTINCT STOR_FLOOR FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :id`, id)
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
	c.JSON(http.StatusOK, floors)
}

// GetOrganizations returns all organizations for a building and floor
func GetOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	floor := c.Param("floor")
	rows, err := database.DB.Query(`SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :id AND STOR_FLOOR = :floor`, id, floor)
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
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng); err != nil {
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
		}
		orgs = append(orgs, orgMap)
	}
	c.JSON(http.StatusOK, orgs)
}

// GetAllOrganizations returns all organizations for a building (all floors)
func GetAllOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rows, err := database.DB.Query(`SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :id`, id)
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
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng); err != nil {
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

	rows, err := database.DB.Query(`SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :id`, id)
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
		}
		if err := rows.Scan(&m.ID, &m.OpTypeName, &m.DistCode, &m.KhoCode, &m.StorName, &m.StorFloor, &m.MrchRegno, &m.PayCenterPropertyID, &m.PayCenterID, &m.Lat, &m.Lng); err != nil {
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
	var totalArea float64

	// 1. Нийт Объект (барилга) - PAY_CENTER тоо
	err := database.DB.QueryRow("SELECT COUNT(*) FROM GPS.PAY_CENTER").Scan(&totalBuildings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting buildings: " + err.Error()})
		return
	}

	// 2. Нийт мкв - PAY_CENTER_PROPERTY доторх бүх PROPERTY_SIZE нэмэх
	err = database.DB.QueryRow("SELECT NVL(SUM(TO_NUMBER(REPLACE(PROPERTY_SIZE, ',', '.'))), 0) FROM GPS.PAY_CENTER_PROPERTY").Scan(&totalArea)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating total area: " + err.Error()})
		return
	}

	// 3. Түрээслэгч - нийт хэдэн түрээслэгч (PAY_MARKET доторх unique MRCH_REGNO)
	err = database.DB.QueryRow("SELECT COUNT(DISTINCT MRCH_REGNO) FROM GPS.PAY_MARKET").Scan(&totalTenants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting tenants: " + err.Error()})
		return
	}

	// 4. Баримт хэвлэдэг - V_E_TUB_PAY_MARKET_EBARIMT доторх бүх мөрийг тоолох
	err = database.DB.QueryRow("SELECT COUNT(*) FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT").Scan(&totalReceiptCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting receipts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_buildings":     totalBuildings,
			"total_area":          totalArea,
			"total_tenants":       totalTenants,
			"total_receipt_count": totalReceiptCount,
		},
	})
}
