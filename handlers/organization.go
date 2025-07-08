package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"dashboard-backend/database"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// ListOrganization returns a list of organizations (GET /api/v1/organization)
// Now returns: ID, NAME, OFFICE_CODE, REGNO, KHO_CODE, BUILD_FLOOR, ADDRESS, LNG, LAT
func ListOrganization(c *gin.Context) {
	name := c.Query("name")
	code := c.Query("code")
	status := c.Query("status")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page_number", "1"))

	orgs, err := repository.GetOrgList(name, code, status, pageSize, pageNumber)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetOrganizationByID returns a single organization by ID (GET /api/v1/organization/:id)
func GetOrganizationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	org, err := repository.FindOrgByID(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if org == nil {
		auth.ErrorResponse(c, http.StatusNotFound, "Organization not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": org})
}

// CreateOrganization creates a new organization (POST /api/v1/organization)
func CreateOrganization(c *gin.Context) {
	var org repository.OrgInput
	if err := c.ShouldBindJSON(&org); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	created, err := repository.CreateOrg(org)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": created})
}

// UpdateOrganization updates an organization (PUT /api/v1/organization/:id)
func UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var org repository.OrgInput
	if err := c.ShouldBindJSON(&org); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated, err := repository.UpdateOrg(uint(id), org)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// DeleteOrganization deletes an organization (DELETE /api/v1/organization/:id)
func DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = repository.DeleteOrg(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}

func ListOrganizations(c *gin.Context) {
	db := database.DB
	rows, err := db.Query("SELECT ID, NAME, REGNO, LNG, LAT FROM GPS.PAY_CENTER")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orgs []map[string]interface{}
	for rows.Next() {
		var id int
		var name, regno string
		var lng, lat sql.NullFloat64

		if err := rows.Scan(&id, &name, &regno, &lng, &lat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		org := map[string]interface{}{
			"id":    id,
			"name":  name,
			"regno": regno,
			"lng":   nil,
			"lat":   nil,
		}

		if lng.Valid {
			org["lng"] = lng.Float64
		}
		if lat.Valid {
			org["lat"] = lat.Float64
		}

		orgs = append(orgs, org)
	}
	c.JSON(http.StatusOK, orgs)
}

// GetOrganizationsBatch returns organizations with all related data in one query
func GetOrganizationsBatch(c *gin.Context) {
	id := c.Query("id")
	floor := c.Query("floor")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	// Use the correct column names from PAY_MARKET table
	payCenterID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var query string
	var args []interface{}

	if floor != "" {
		query = `SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :1 AND STOR_FLOOR = :2`
		args = []interface{}{payCenterID, floor}
	} else {
		query = `SELECT ID, OP_TYPE_NAME, DIST_CODE, KHO_CODE, STOR_NAME, STOR_FLOOR, MRCH_REGNO, PAY_CENTER_PROPERTY_ID, PAY_CENTER_ID, LAT, LNG FROM GPS.PAY_MARKET WHERE PAY_CENTER_ID = :1`
		args = []interface{}{payCenterID}
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
		return
	}
	defer rows.Close()

	var orgs []map[string]interface{}
	for rows.Next() {
		var (
			id                  int
			opTypeName          sql.NullString
			distCode            sql.NullString
			khoCode             sql.NullString
			storName            sql.NullString
			storFloor           sql.NullString
			mrchRegno           sql.NullString
			payCenterPropertyID sql.NullInt64
			payCenterID         int
			lat                 sql.NullFloat64
			lng                 sql.NullFloat64
		)

		if err := rows.Scan(&id, &opTypeName, &distCode, &khoCode, &storName, &storFloor, &mrchRegno, &payCenterPropertyID, &payCenterID, &lat, &lng); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
			return
		}

		org := map[string]interface{}{
			"id":                     id,
			"op_type_name":           getStringValue(opTypeName),
			"dist_code":              getStringValue(distCode),
			"kho_code":               getStringValue(khoCode),
			"stor_name":              getStringValue(storName),
			"stor_floor":             getStringValue(storFloor),
			"mrch_regno":             getStringValue(mrchRegno),
			"pay_center_property_id": getInt64Value(payCenterPropertyID),
			"pay_center_id":          payCenterID,
			"lat":                    getFloatValue(lat),
			"lng":                    getFloatValue(lng),
			"count_receipt":          0,  // Default for now
			"report_submitted_date":  "", // Default for now
			"payable_debit":          0,  // Default for now
			"advice_count":           0,  // Default for now
		}

		orgs = append(orgs, org)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// Helper functions to handle NULL values
func getStringValue(nullStr sql.NullString) interface{} {
	if nullStr.Valid {
		return nullStr.String
	}
	return nil
}

func getFloatValue(nullFloat sql.NullFloat64) interface{} {
	if nullFloat.Valid {
		return nullFloat.Float64
	}
	return nil
}

// Helper function for Int64 values
func getInt64Value(nullInt sql.NullInt64) interface{} {
	if nullInt.Valid {
		return int(nullInt.Int64)
	}
	return nil
}
