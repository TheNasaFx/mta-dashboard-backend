package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/database"
	"dashboard-backend/repository"
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
	var results []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		BuildFloor int    `json:"build_floor"`
		OfficeCode string `json:"office_code"`
		KhoCode    string `json:"kho_code"`
		Regno      string `json:"regno"`
		Lat        string `json:"lat"`
		Lng        string `json:"lng"`
		ParcelId   string `json:"parcel_id"`
	}
	for rows.Next() {
		var pc struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			BuildFloor int    `json:"build_floor"`
			OfficeCode string `json:"office_code"`
			KhoCode    string `json:"kho_code"`
			Regno      string `json:"regno"`
			Lat        string `json:"lat"`
			Lng        string `json:"lng"`
			ParcelId   string `json:"parcel_id"`
		}
		if err := rows.Scan(&pc.ID, &pc.Name, &pc.BuildFloor, &pc.OfficeCode, &pc.KhoCode, &pc.Regno, &pc.Lat, &pc.Lng, &pc.ParcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		results = append(results, pc)
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
	var results []struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		BuildFloor int    `json:"build_floor"`
		OfficeCode string `json:"office_code"`
		KhoCode    string `json:"kho_code"`
		Regno      string `json:"regno"`
		Lat        string `json:"lat"`
		Lng        string `json:"lng"`
		ParcelId   string `json:"parcel_id"`
	}
	for rows.Next() {
		var pc struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			BuildFloor int    `json:"build_floor"`
			OfficeCode string `json:"office_code"`
			KhoCode    string `json:"kho_code"`
			Regno      string `json:"regno"`
			Lat        string `json:"lat"`
			Lng        string `json:"lng"`
			ParcelId   string `json:"parcel_id"`
		}
		if err := rows.Scan(&pc.ID, &pc.Name, &pc.BuildFloor, &pc.OfficeCode, &pc.KhoCode, &pc.Regno, &pc.Lat, &pc.Lng, &pc.ParcelId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}
		results = append(results, pc)
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
