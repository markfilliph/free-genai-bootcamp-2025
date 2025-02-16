package models

import (
	"database/sql"
	"fmt"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB initializes the database connection
func InitDB(dbPath string) error {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return
		}

		// Set connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)

		// Test the connection
		if err = db.Ping(); err != nil {
			return
		}

		// Enable foreign key constraints
		if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			return
		}
	})

	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// Transaction represents a database transaction
type Transaction struct {
	*sql.Tx
}

// Begin starts a new transaction
func Begin() (*Transaction, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return &Transaction{tx}, nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// Rollback rolls back the transaction
func (tx *Transaction) Rollback() error {
	if err := tx.Tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %v", err)
	}
	return nil
}
