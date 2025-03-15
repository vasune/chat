package database

import (
	"chat/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func SetupDB() {
	config.LoadCfg()
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		config.AppConfig.PostgresUser,
		config.AppConfig.PostgresPassword,
		config.AppConfig.PostgresHost,
		config.AppConfig.PostgresPort,
		config.AppConfig.PostgresDB,
	)

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
			)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database initialized")
}
