package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"backend/database"
)

func main() {
	// DB холболт
	database.MustConnect()

	// health check endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		var version string
		err := database.DB.QueryRowContext(context.Background(), "SELECT banner FROM v$version WHERE ROWNUM = 1").Scan(&version)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB алдаа: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "🟢 Oracle DB холбогдсон: %s", version)
	})

	log.Println("Сервер http://localhost:8080 дээр ажиллаж байна...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
