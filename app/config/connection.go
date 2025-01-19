package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
)

var errEnv = godotenv.Load()
var AppPort = os.Getenv("APP_PORT")
var UserAdmin = os.Getenv("USER_ADMIN")
var PassAdmin = os.Getenv("PASS_ADMIN")

func ConnectToDatabase() (*sql.DB, error) {
	if errEnv != nil {
		return nil, fmt.Errorf("failed to load file .env: %v", errEnv)
	}

	db, err := sql.Open("sqlite3", "./pendaftar.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while connecting to the database: %v", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}