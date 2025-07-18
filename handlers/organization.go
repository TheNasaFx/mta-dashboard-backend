package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/database"
	"dashboard-backend/repository"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
	rows, err := db.Query("SELECT ID, NAME, REGNO, LNG, LAT, BUILD_FLOOR FROM GPS.PAY_CENTER")
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
		var buildFloor sql.NullInt64

		if err := rows.Scan(&id, &name, &regno, &lng, &lat, &buildFloor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		org := map[string]interface{}{
			"id":          id,
			"name":        name,
			"regno":       regno,
			"lng":         nil,
			"lat":         nil,
			"build_floor": nil,
		}

		if lng.Valid {
			org["lng"] = lng.Float64
		}
		if lat.Valid {
			org["lat"] = lat.Float64
		}
		if buildFloor.Valid {
			org["build_floor"] = int(buildFloor.Int64)
		}

		orgs = append(orgs, org)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
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
		query = `SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR FROM GPS.PAY_MARKET pm LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID WHERE pm.PAY_CENTER_ID = :1 AND pm.STOR_FLOOR = :2`
		args = []interface{}{payCenterID, floor}
	} else {
		query = `SELECT pm.ID, pm.OP_TYPE_NAME, pm.DIST_CODE, pm.KHO_CODE, pm.STOR_NAME, pm.STOR_FLOOR, pm.MRCH_REGNO, pm.PAY_CENTER_PROPERTY_ID, pm.PAY_CENTER_ID, pm.LAT, pm.LNG, pc.BUILD_FLOOR FROM GPS.PAY_MARKET pm LEFT JOIN GPS.PAY_CENTER pc ON pm.PAY_CENTER_ID = pc.ID WHERE pm.PAY_CENTER_ID = :1`
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
			buildFloor          sql.NullInt64
		)

		if err := rows.Scan(&id, &opTypeName, &distCode, &khoCode, &storName, &storFloor, &mrchRegno, &payCenterPropertyID, &payCenterID, &lat, &lng, &buildFloor); err != nil {
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
			"build_floor":            getInt64Value(buildFloor),
			"count_receipt":          0,  // Default for now
			"report_submitted_date":  "", // Default for now
			"payable_debit":          0,  // Default for now
			"advice_count":           0,  // Default for now
		}

		orgs = append(orgs, org)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetOrganizationDetail returns detailed information about an organization
func GetOrganizationDetail(c *gin.Context) {
	mrchRegno := c.Param("regno")
	if mrchRegno == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "regno parameter is required"})
		return
	}

	result := gin.H{
		"success": true,
		"data": gin.H{
			"mrch_regno": mrchRegno,
			"branch":     nil,
			"segment":    nil,
			"ebarimt": gin.H{
				"cnt_3":  0,
				"cnt_30": 0,
			},
			"reports":          []interface{}{},
			"payments":         []interface{}{},
			"debts":            []interface{}{},
			"license_info":     0,
			"advisory_service": 0,
			"violation_info":   0,
		},
	}

	// 1. V_E_TUB_BRANCH мэдээлэл авах
	branchQuery := `SELECT REGISTER, OVOG_NER, TTA, DED_ALBA, TULUV 
		FROM GPS.V_E_TUB_BRANCH 
		WHERE TRIM(UPPER(REGISTER)) = TRIM(UPPER(:1))`

	var branch struct {
		Register sql.NullString
		OvogNer  sql.NullString
		TTA      sql.NullString
		DedAlba  sql.NullString
		Tuluv    sql.NullString
	}

	err := database.DB.QueryRow(branchQuery, mrchRegno).Scan(
		&branch.Register, &branch.OvogNer, &branch.TTA, &branch.DedAlba, &branch.Tuluv)

	if err == nil {
		branchData := gin.H{}
		if branch.Register.Valid {
			branchData["register"] = branch.Register.String
		}
		if branch.OvogNer.Valid {
			branchData["ovog_ner"] = branch.OvogNer.String
		}
		if branch.TTA.Valid {
			branchData["tta"] = branch.TTA.String
		}
		if branch.DedAlba.Valid {
			branchData["ded_alba"] = branch.DedAlba.String
		}
		if branch.Tuluv.Valid {
			branchData["tuluv"] = branch.Tuluv.String
		}
		result["data"].(gin.H)["branch"] = branchData
	}

	// 2. V_E_TUB_SEGMENT мэдээлэл авах
	segmentQuery := `SELECT SEGMENT, SEGMENT_YEAR 
		FROM GPS.V_E_TUB_SEGMENT 
		WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))`

	var segment struct {
		Segment     sql.NullString
		SegmentYear sql.NullString
	}

	err = database.DB.QueryRow(segmentQuery, mrchRegno).Scan(&segment.Segment, &segment.SegmentYear)
	if err == nil {
		segmentData := gin.H{}
		if segment.Segment.Valid {
			segmentData["segment"] = segment.Segment.String
		}
		if segment.SegmentYear.Valid {
			segmentData["segment_year"] = segment.SegmentYear.String
		}
		result["data"].(gin.H)["segment"] = segmentData
	}

	// 3. Е-баримт мэдээлэл (өмнө хийсэн API-аас)
	ebarimtQuery := `SELECT 
		COALESCE(SUM(CNT_3), 0) as CNT_3,
		COALESCE(SUM(CNT_30), 0) as CNT_30
		FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT 
		WHERE TRIM(UPPER(MRCH_REGNO)) = TRIM(UPPER(:1))`

	var cnt3, cnt30 int
	err = database.DB.QueryRow(ebarimtQuery, mrchRegno).Scan(&cnt3, &cnt30)
	if err == nil {
		result["data"].(gin.H)["ebarimt"] = gin.H{
			"cnt_3":  cnt3,
			"cnt_30": cnt30,
		}
	}

	// 4. V_E_TUB_REPORT_DATA тайлангийн мэдээлэл
	reportQuery := `SELECT TAX_REPORT_CODE, FREQUENCY, TAX_YEAR, TAX_PERIOD, SUBMITTED_DATE 
		FROM GPS.V_E_TUB_REPORT_DATA 
		WHERE TRIM(UPPER(TIN)) = TRIM(UPPER(:1))
		ORDER BY SUBMITTED_DATE DESC`

	rows, err := database.DB.Query(reportQuery, mrchRegno)
	if err == nil {
		defer rows.Close()
		reports := []gin.H{}
		for rows.Next() {
			var taxReportCode, frequency, taxYear, taxPeriod, submittedDate sql.NullString
			if err := rows.Scan(&taxReportCode, &frequency, &taxYear, &taxPeriod, &submittedDate); err == nil {
				report := gin.H{}
				if taxReportCode.Valid {
					report["tax_report_code"] = taxReportCode.String
				}
				if frequency.Valid {
					report["frequency"] = frequency.String
				}
				if taxYear.Valid {
					report["tax_year"] = taxYear.String
				}
				if taxPeriod.Valid {
					report["tax_period"] = taxPeriod.String
				}
				if submittedDate.Valid {
					report["submitted_date"] = submittedDate.String
				}
				reports = append(reports, report)
			}
		}
		result["data"].(gin.H)["reports"] = reports
	}

	// 5. V_E_TUB_PAYMENTS төлөлтийн мэдээлэл (PIN талбараар хайх)
	paymentQuery := `SELECT INVOICE_NO, TAX_TYPE_NAME, BRANCH_NAME, AMOUNT, PAID_DATE 
		FROM GPS.V_E_TUB_PAYMENTS 
		WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1))
		ORDER BY PAID_DATE DESC`

	rows, err = database.DB.Query(paymentQuery, mrchRegno)
	if err != nil {
		fmt.Printf("Payment query error: %v\n", err)
	} else {
		defer rows.Close()
		payments := []gin.H{}
		for rows.Next() {
			var invoiceNo, taxTypeName, branchName sql.NullString
			var amount sql.NullFloat64
			var paidDate sql.NullString
			if err := rows.Scan(&invoiceNo, &taxTypeName, &branchName, &amount, &paidDate); err == nil {
				payment := gin.H{}
				if invoiceNo.Valid {
					payment["invoice_no"] = invoiceNo.String
				}
				if taxTypeName.Valid {
					payment["tax_type_name"] = taxTypeName.String
				}
				if branchName.Valid {
					payment["branch_name"] = branchName.String
				}
				if amount.Valid {
					payment["amount"] = amount.Float64
				}
				if paidDate.Valid {
					payment["paid_date"] = paidDate.String
				}
				payments = append(payments, payment)
			}
		}
		fmt.Printf("Found %d payments for regno %s\n", len(payments), mrchRegno)
		result["data"].(gin.H)["payments"] = payments
	}

	// 6. V_ACCOUNT_GENERAL_YEAR өрийн мэдээлэл
	debtQuery := `SELECT TAX_TYPE_NAME, YEAR, PERIOD_TYPE, BRANCH_NAME, C2_DEBIT 
		FROM GPS.V_ACCOUNT_GENERAL_YEAR 
		WHERE TRIM(UPPER(PIN)) = TRIM(UPPER(:1)) AND C2_DEBIT > 0
		ORDER BY YEAR DESC, PERIOD_TYPE`

	rows, err = database.DB.Query(debtQuery, mrchRegno)
	if err == nil {
		defer rows.Close()
		debts := []gin.H{}
		for rows.Next() {
			var taxTypeName, year, periodType, branchName sql.NullString
			var c2Debit sql.NullFloat64
			if err := rows.Scan(&taxTypeName, &year, &periodType, &branchName, &c2Debit); err == nil {
				debt := gin.H{}
				if taxTypeName.Valid {
					debt["tax_type_name"] = taxTypeName.String
				}
				if year.Valid {
					debt["year"] = year.String
				}
				if periodType.Valid {
					debt["period_type"] = periodType.String
				}
				if branchName.Valid {
					debt["branch_name"] = branchName.String
				}
				if c2Debit.Valid {
					debt["c2_debit"] = c2Debit.Float64
				}
				debts = append(debts, debt)
			}
		}
		result["data"].(gin.H)["debts"] = debts
	}

	c.JSON(http.StatusOK, result)
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
