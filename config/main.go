package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	Service  string
	User     string
	Password string
}

type Config struct {
	DB DBConfig
}

func Get() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env файл олдсонгүй. Шууд os.Getenv ашиглана.")
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			Service:  os.Getenv("DB_SERVICE"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}
