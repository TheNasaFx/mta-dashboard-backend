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
	c.JSON(http.StatusOK, orgs)
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

// GetEbarimtByPin returns ebarimt info by PIN
func GetEbarimtByPin(c *gin.Context) {
	pin := c.Param("pin")
	ebarimt, err := repository.GetEbarimtByPin(pin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ebarimt not found for this PIN"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ebarimt})
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

	// Check if this MRCH_REGNO exists as PIN in PAY_MARKET_BARIMT
	var barimtExists int
	var countReceipt int
	err = database.DB.QueryRow("SELECT COUNT(*), COALESCE(MAX(COUNT_RECEIPT), 0) FROM GPS.PAY_MARKET_BARIMT WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))", mrchRegno).Scan(&barimtExists, &countReceipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking PAY_MARKET_BARIMT: " + err.Error()})
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
