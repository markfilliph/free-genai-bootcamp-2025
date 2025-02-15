package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// GetDB returns the database connection
func GetDB() (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = InitDB()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
		db = nil
	}
	return nil
}

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
