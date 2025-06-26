package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"dashboard-backend/database"
	"dashboard-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
func getBuildings(c *gin.Context) {
	c.JSON(http.StatusOK, centers)
}

func getCenters(c *gin.Context) {
	district := c.Query("district")
	khoroo := c.Query("khoroo")
	filtered := []PayCenter{}
	for _, center := range centers {
		if (district == "" || center.OfficeCode == district) && (khoroo == "" || center.KhoCode == khoroo) {
			filtered = append(filtered, center)
		}
	}
	marketMrchMap := make(map[int][]string)
	for _, m := range markets {
		marketMrchMap[m.PayCenterID] = append(marketMrchMap[m.PayCenterID], m.MrchRegno)
	}
	barimtData, err := ioutil.ReadFile("data/pay_market_barimt.json")
	type BarimtItem struct {
		ID           int    `json:"id"`
		Pin          string `json:"pin"`
		CountReceipt int    `json:"count_receipt"`
	}
	var barimtRaw struct {
		Results []struct{ Items []BarimtItem }
	}
	barimtMap := make(map[string]int)
	if err == nil {
		json.Unmarshal(barimtData, &barimtRaw)
		for _, item := range barimtRaw.Results[0].Items {
			barimtMap[item.Pin] = item.CountReceipt
		}
	}
	type CenterWithPin struct {
		PayCenter
		PinList     []string `json:"pin_list"`
		AllBarimtOk bool     `json:"all_barimt_ok"`
	}
	result := []CenterWithPin{}
	for _, center := range filtered {
		pins := marketMrchMap[center.ID]
		allOk := true
		for _, pin := range pins {
			if barimtMap[pin] == 0 {
				allOk = false
				break
			}
		}
		result = append(result, CenterWithPin{
			PayCenter:   center,
			PinList:     pins,
			AllBarimtOk: allOk,
		})
	}
	c.JSON(http.StatusOK, result)
}

func getFloors(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	floorSet := make(map[string]bool)
	for _, m := range markets {
		if m.PayCenterID == id {
			floorSet[m.StorFloor] = true
		}
	}
	floors := []string{}
	for f := range floorSet {
		floors = append(floors, f)
	}
	c.JSON(http.StatusOK, floors)
}

func getOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	floor := c.Param("floor")
	orgs := []map[string]interface{}{}
	barimtData, err := ioutil.ReadFile("data/pay_market_barimt.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pay_market_barimt.json not found"})
		return
	}
	var barimtRaw struct {
		Results []struct {
			Items []struct {
				ID           int    `json:"id"`
				Pin          string `json:"pin"`
				CountReceipt int    `json:"count_receipt"`
			}
		}
	}
	err = json.Unmarshal(barimtData, &barimtRaw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pay_market_barimt.json parse error"})
		return
	}
	barimtMap := make(map[string]int)
	if len(barimtRaw.Results) > 0 {
		for _, item := range barimtRaw.Results[0].Items {
			barimtMap[item.Pin] = item.CountReceipt
		}
	}
	for _, m := range markets {
		if m.PayCenterID == id && m.StorFloor == floor {
			org := m
			orgMap := make(map[string]interface{})
			b, _ := json.Marshal(org)
			json.Unmarshal(b, &orgMap)
			orgMap["count_receipt"] = barimtMap[m.MrchRegno]
			orgMap["lat"] = m.Lat
			orgMap["lng"] = m.Lng
			orgs = append(orgs, orgMap)
		}
	}
	c.JSON(http.StatusOK, orgs)
}

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

func getAllOrganizations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	orgs := []map[string]interface{}{}
	for _, m := range markets {
		if m.PayCenterID == id {
			org := m
			orgMap := make(map[string]interface{})
			b, err := json.Marshal(org)
			if err != nil {
				continue
			}
			json.Unmarshal(b, &orgMap)
			orgs = append(orgs, orgMap)
		}
	}
	c.JSON(http.StatusOK, orgs)
}

func main() {
	// pay_center.json унших
	centerData, err := ioutil.ReadFile("data/pay_center.json")
	if err != nil {
		panic(err)
	}
	var centerRaw struct{ Results []struct{ Items []PayCenter } }
	err = json.Unmarshal(centerData, &centerRaw)
	if err != nil {
		panic(err)
	}
	if len(centerRaw.Results) > 0 {
		centers = centerRaw.Results[0].Items
	}

	// pay_market.json унших
	marketData, err := ioutil.ReadFile("data/pay_market.json")
	if err != nil {
		panic(err)
	}
	var marketRaw struct{ Results []struct{ Items []PayMarket } }
	err = json.Unmarshal(marketData, &marketRaw)
	if err != nil {
		panic(err)
	}
	if len(marketRaw.Results) > 0 {
		markets = marketRaw.Results[0].Items
	}

	database.MustConnect()

	router := gin.Default()

	// CORS middleware нэмэх
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/api/buildings", getBuildings)
	router.GET("/api/centers", getCenters)
	router.GET("/api/buildings/:id/floors", getFloors)
	router.GET("/api/buildings/:id/floors/:floor/organizations", getOrganizations)
	router.GET("/api/organizations/:id", getOrganizationDetail)
	router.GET("/api/barimt", func(c *gin.Context) {
		w := c.Writer
		w.Header().Set("Content-Type", "application/json")
		data, err := os.ReadFile("data/pay_market_barimt.json")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Write(data)
	})
	router.GET("/api/buildings/:id/organizations", getAllOrganizations)

	routes.RegisterV1Routes(router)
	router.Run(":8080")
}
