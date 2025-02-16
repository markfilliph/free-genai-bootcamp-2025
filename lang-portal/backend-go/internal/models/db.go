package models

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"log"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() (*sql.DB, error) {
	if DB == nil {
		return nil, sql.ErrConnDone
	}
	return DB, nil
}