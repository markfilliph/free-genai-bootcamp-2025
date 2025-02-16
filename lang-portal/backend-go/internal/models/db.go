package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the SQLite database connection
func InitDB() error {
	var err error

	// Database file will be in the root of the project
	dbPath := filepath.Join(".", "words.db")
	log.Printf("Opening SQLite database at: %s", dbPath)

	// Open SQLite database
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("error enabling foreign keys: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// RunMigrations runs the database migrations from the specified directory
func RunMigrations() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Read migration file
	migrationPath := filepath.Join("db", "migrations", "001_initial_schema.sql")
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migration file: %v", err)
	}

	// Execute migration
	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("error executing migration: %v", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// RunSeeds runs the database seeds from the specified directory
func RunSeeds() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Read seed file
	seedPath := filepath.Join("db", "seeds", "initial_data.sql")
	seed, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("error reading seed file: %v", err)
	}

	// Execute seed
	_, err = db.Exec(string(seed))
	if err != nil {
		return fmt.Errorf("error executing seed: %v", err)
	}

	log.Println("Database seeding completed")
	return nil
}
