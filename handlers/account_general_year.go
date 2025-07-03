package handlers

import (
	"dashboard-backend/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccountGeneralYearsHandler(c *gin.Context) {
	regno := c.Query("regno")
	db := database.DB
	tab := c.Query("tab") // 'report' or 'info'
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
			c.JSON(http.StatusOK, results)
			return
		} else if tab == "info" {
			rows, err := db.Query("SELECT PIN, ENTITY_NAME FROM GPS.V_ACCOUNT_GENERAL_YEAR WHERE PIN = :1", regno)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			type Info struct {
				Pin        string `json:"pin"`
				EntityName string `json:"entity_name"`
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
			c.JSON(http.StatusOK, results)
			return
		}
	}
	c.JSON(http.StatusOK, []interface{}{})
}
