package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppConfig Appcfg

type Appcfg struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	JwtSecretKey     string
}

func LoadCfg() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error reading .env: %s", err)
		return
	}
	AppConfig = Appcfg{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		JwtSecretKey:     os.Getenv("JWT_SECRET"),
	}

}
