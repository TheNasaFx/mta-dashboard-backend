package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"dashboard-backend/config"

	_ "github.com/sijms/go-ora/v2"
)

var DB *sql.DB

func MustConnect() {
	cfg := config.Get()

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Service,
	)

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		log.Fatalf("Oracle өгөгдлийн сантай холбогдох үед алдаа гарлаа: %v", err)
	}

	// Connection pool тохиргоо
	db.SetMaxOpenConns(25)                  // Хамгийн их 25 холболт
	db.SetMaxIdleConns(5)                   // Хамгийн их 5 idle холболт
	db.SetConnMaxLifetime(5 * time.Minute)  // Connection-ий max хугацаа
	db.SetConnMaxIdleTime(30 * time.Second) // Idle connection timeout

	if err = db.Ping(); err != nil {
		log.Fatalf("Oracle өгөгдлийн сан ping амжилтгүй: %v", err)
	}

	log.Println("Oracle өгөгдлийн сантай амжилттай холбогдлоо.")
	log.Printf("Connection pool: MaxOpen=%d, MaxIdle=%d", 25, 5)
	DB = db
}
