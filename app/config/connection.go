package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func ConnectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "database/pendaftar.db")
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