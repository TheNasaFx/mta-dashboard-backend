package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"dashboard-backend/database"
	"dashboard-backend/handlers"
	"dashboard-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go backend!")
}

func checkOracleConnection(w http.ResponseWriter, r *http.Request) {
	if database.DB == nil {
		database.MustConnect()
	}

	err := database.DB.Ping()
	if err != nil {
		http.Error(w, "❌ Oracle DB connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "✅ Oracle DB connection is healthy.")
}

type PayCenter struct {
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

type PayMarket struct {
	ID          int    `json:"id"`
	OpTypeName  string `json:"op_type_name"`
	DistCode    string `json:"dist_code"`
	KhoCode     string `json:"kho_code"`
	MarCode     string `json:"mar_code"`
	MarName     string `json:"mar_name"`
	MarRegno    string `json:"mar_regno"`
	StorCode    string `json:"stor_code"`
	StorFloor   string `json:"stor_floor"`
	StorName    string `json:"stor_name"`
	PayCenterID int    `json:"pay_center_id"`
	MrchRegno   string `json:"mrch_regno"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
}

var centers []PayCenter
var markets []PayMarket

// Login handler struct and logic
// (from handlers/login.go)
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	ActiveValue int    `json:"activeValue"`
	MonitorId   string `json:"monitorId"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}
	externalURL := "https://st-tais.mta.mn/rest/tais-ims-service/token/login"
	externalReq, err := http.NewRequest("POST", externalURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create external request"})
		return
	}
	externalReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	externalResp, err := client.Do(externalReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to external API"})
		return
	}
	defer externalResp.Body.Close()
	body, err := io.ReadAll(externalResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read external response"})
		return
	}
	if externalResp.StatusCode != http.StatusOK {
		c.JSON(externalResp.StatusCode, gin.H{"error": string(body)})
		return
	}
	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse external response"})
		return
	}
	c.JSON(http.StatusOK, loginResp)
}

func getPayCenters(w http.ResponseWriter, r *http.Request) {
	if database.DB == nil {
		database.MustConnect()
	}

	rows, err := database.DB.Query(`SELECT ID, NAME FROM GPS.PAY_CENTER`)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []PayCenter
	for rows.Next() {
		var pc PayCenter
		if err := rows.Scan(&pc.ID, &pc.Name); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, pc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// --- Data API handler functions ---
// getFloors, getOrganizations, getAllOrganizations функцуудыг устгана. Одоо зөвхөн handlers/market.go-оос импортолно.

func getOrganizationDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, m := range markets {
		if m.ID == id {
			c.JSON(http.StatusOK, m)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func main() {
	_ = godotenv.Load()
	database.MustConnect()
	secret := os.Getenv("JWT_SECRET")
	fmt.Println("JWT_SECRET:", secret) // Түр зуур хэвлэж шалга

	router := gin.Default()

	// CORS middleware нэмэх
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/api/organizations/:id", getOrganizationDetail)
	router.GET("/api/barimt", func(c *gin.Context) {
		rows, err := database.DB.Query(`SELECT ID, MRCH_REGNO, CNT_3, CNT_30, OP_TYPE_NAME, MAR_NAME, MAR_REGNO, QR_CODE FROM GPS.V_E_TUB_PAY_MARKET_EBARIMT`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
			return
		}
		defer rows.Close()
		var results []map[string]interface{}
		for rows.Next() {
			var id int
			var mrchRegno, opTypeName, marName, marRegno, qrCode sql.NullString
			var cnt3, cnt30 sql.NullInt64
			if err := rows.Scan(&id, &mrchRegno, &cnt3, &cnt30, &opTypeName, &marName, &marRegno, &qrCode); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
				return
			}
			results = append(results, map[string]interface{}{
				"id":            id,
				"mrch_regno":    ifNullString(mrchRegno),
				"cnt_3":         ifNullInt64(cnt3),
				"cnt_30":        ifNullInt64(cnt30),
				"count_receipt": ifNullInt64(cnt3), // Backward compatibility with existing frontend
				"op_type_name":  ifNullString(opTypeName),
				"mar_name":      ifNullString(marName),
				"mar_regno":     ifNullString(marRegno),
				"qr_code":       ifNullString(qrCode),
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
	})
	router.GET("/api/pay_center_property", func(c *gin.Context) {
		rows, err := database.DB.Query(`SELECT ID, PAY_CENTER_ID, PROPERTY_NUMBER, STATUS, CREATED_BY, CREATED_DATE, UPDATED_BY, UPDATED_DATE, PROPERTY_TYPE, OWNER_REGNO, PROPERTY_SIZE, RENT_AMOUNT FROM GPS.PAY_CENTER_PROPERTY`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query error: " + err.Error()})
			return
		}
		defer rows.Close()
		var results []map[string]interface{}
		for rows.Next() {
			var id, payCenterId sql.NullInt64
			var propertyNumber, ownerRegno, createdDate, updatedDate sql.NullString
			var status, createdBy, updatedBy, propertyType sql.NullInt64
			var propertySize, rentAmount sql.NullFloat64
			if err := rows.Scan(&id, &payCenterId, &propertyNumber, &status, &createdBy, &createdDate, &updatedBy, &updatedDate, &propertyType, &ownerRegno, &propertySize, &rentAmount); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
				return
			}
			results = append(results, map[string]interface{}{
				"id":              ifNullInt64(id),
				"pay_center_id":   ifNullInt64(payCenterId),
				"property_number": ifNullString(propertyNumber),
				"status":          ifNullInt64(status),
				"created_by":      ifNullInt64(createdBy),
				"created_date":    ifNullString(createdDate),
				"updated_by":      ifNullInt64(updatedBy),
				"updated_date":    ifNullString(updatedDate),
				"property_type":   ifNullInt64(propertyType),
				"owner_regno":     ifNullString(ownerRegno),
				"property_size":   ifNullFloat64(propertySize),
				"rent_amount":     ifNullFloat64(rentAmount),
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
	})
	router.GET("/api/ebarimt/:pin", handlers.GetEbarimtByPin)
	router.GET("/api/buildings/:id/floors", handlers.GetFloors)
	router.GET("/api/buildings/:id/floors/:floor/organizations", handlers.GetOrganizations)
	router.GET("/api/buildings/:id/organizations", handlers.GetAllOrganizations)

	routes.RegisterV1Routes(router)
	router.Run(":8080")
}

// Helper functions for handling NULL values
func ifNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ifNullFloat64(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}

func ifNullInt64(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return 0
}
