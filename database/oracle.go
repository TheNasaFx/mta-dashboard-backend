package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend/config"

	_ "github.com/sijms/go-ora/v2"
)

var DB *sql.DB

func MustConnect() {
	cfg := config.Get()

	// Oracle DSN хэлбэр:
	// oracle://user:password@host:port/service_name
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

	// Холболтыг шалгах
	if err = db.Ping(); err != nil {
		log.Fatalf("Oracle өгөгдлийн сан ping амжилтгүй: %v", err)
	}

	log.Println("✅ Oracle өгөгдлийн сантай амжилттай холбогдлоо.")
	DB = db
}
