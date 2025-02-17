package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize sets up the database connection and runs migrations
func Initialize() error {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Open SQLite database
	dbPath := filepath.Join(wd, "words.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db

	// Run migrations
	if err := runMigrations(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// runMigrations executes all migration files in order
func runMigrations() error {
	// Read migration files
	migrationsDir := "./db/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	// Execute each migration file
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return err
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return err
		}

		log.Printf("Executed migration: %s\n", file.Name())
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
