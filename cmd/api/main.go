package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dashboard-backend/database"
	"dashboard-backend/routes"

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
	ID   int    `json:"id"`
	Name string `json:"name"`
	// ...таны хүснэгтийн бусад багануудыг энд нэмнэ үү...
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

func main() {
	database.MustConnect()

	router := gin.Default()
	routes.RegisterV1Routes(router)
	router.Run(":8080")
}
