package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// DB wraps sql.DB to add custom functionality
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(config *Config) (*DB, error) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
