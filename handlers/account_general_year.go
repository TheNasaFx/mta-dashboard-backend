package handlers

import (
	"dashboard-backend/database"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccountGeneralYearsHandler(c *gin.Context) {
	regno := c.Query("regno")
	db := database.DB
	tab := c.Query("tab") // 'report' or 'info' or 'debt'
	if regno != "" {
		if tab == "report" {
			pageStr := c.DefaultQuery("page", "1")
			sizeStr := c.DefaultQuery("size", "20")
			page, _ := strconv.Atoi(pageStr)
			size, _ := strconv.Atoi(sizeStr)
			if page < 1 {
				page = 1
			}
			if size < 1 {
				size = 20
			}
			offset := (page - 1) * size
			query := `SELECT TAX_TYPE_NAME, TAX_TYPE_CODE, BRANCH_NAME FROM GPS.V_ACCOUNT_GENERAL_YEAR WHERE PIN = :1 OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`
			rows, err := db.Query(query, regno, offset, size)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			type Report struct {
				TaxTypeName string `json:"tax_type_name"`
				TaxTypeCode string `json:"tax_type_code"`
				BranchName  string `json:"branch_name"`
			}
			var results []Report
			for rows.Next() {
				var r Report
				if err := rows.Scan(&r.TaxTypeName, &r.TaxTypeCode, &r.BranchName); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				results = append(results, r)
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
			return
		} else if tab == "info" {
			// Account general year мэдээлэл - зөвхөн нэг мөр буцаана
			rows, err := db.Query("SELECT DISTINCT PIN, ENTITY_NAME FROM GPS.V_ACCOUNT_GENERAL_YEAR WHERE PIN = :1", regno)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			type Info struct {
				Pin          string `json:"pin"`
				EntityName   string `json:"entity_name"`
				EbarimtCount int    `json:"ebarimt_count"`
			}
			var results []Info
			for rows.Next() {
				var i Info
				if err := rows.Scan(&i.Pin, &i.EntityName); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				results = append(results, i)
			}

			// Ebarimt тоог авах
			var ebarimtCount int
			err = db.QueryRow("SELECT COALESCE(MAX(CNT_3), 0) FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT WHERE TRIM(UPPER(MRCH_REGNO)) = TRIM(UPPER(:1))", regno).Scan(&ebarimtCount)
			if err != nil {
				// Алдаа гарвал 0 болгоно
				ebarimtCount = 0
			}

			// Ebarimt тоог бүх мөрөнд нэмэх
			for i := range results {
				results[i].EbarimtCount = ebarimtCount
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
			return
		} else if tab == "debt" {
			rows, err := db.Query("SELECT C2_DEBIT FROM GPS.V_ACCOUNT_GENERAL_YEAR WHERE PIN = :1", regno)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			var results []string
			for rows.Next() {
				var payable sql.NullString
				if err := rows.Scan(&payable); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				if payable.Valid {
					results = append(results, payable.String)
				}
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
}
