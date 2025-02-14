//go:build mage
package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile     = "words.db"
	backupDir  = "db/backups"
	schemaFile = "db/migrations/001_initial_schema.sql"
	seedFile   = "db/seeds/initial_data.sql"
)

// Install installs project dependencies
func Install() error {
	fmt.Println("Installing dependencies...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

// InitDB initializes the database with schema and seed data
func InitDB() error {
	fmt.Println("Initializing database...")

	if err := ensureDBDirectory(); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Create database if it doesn't exist
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Read and execute schema
	if err := executeSQLFile(db, schemaFile); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	// Read and execute seed data
	if err := executeSQLFile(db, seedFile); err != nil {
		return fmt.Errorf("failed to execute seed data: %v", err)
	}

	fmt.Println("Database initialized successfully")
	return nil
}

// Backup creates a backup of the database
func Backup() error {
	if err := ensureBackupDirectory(); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Check if database exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist")
	}

	// Create backup file with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("words_%s.db", timestamp))

	// Copy database file
	if err := copyFile(dbFile, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	fmt.Printf("Database backup created: %s\n", backupFile)
	return nil
}

// Restore restores the database from the most recent backup
func Restore() error {
	// Find most recent backup
	backups, err := filepath.Glob(filepath.Join(backupDir, "words_*.db"))
	if err != nil {
		return fmt.Errorf("failed to list backups: %v", err)
	}
	if len(backups) == 0 {
		return fmt.Errorf("no backups found")
	}

	// Get most recent backup (last file alphabetically, due to timestamp format)
	mostRecent := backups[len(backups)-1]

	// Stop server if running (implement this based on your process management)
	Clean()

	// Restore backup
	if err := copyFile(mostRecent, dbFile); err != nil {
		return fmt.Errorf("failed to restore backup: %v", err)
	}

	fmt.Printf("Database restored from: %s\n", mostRecent)
	return nil
}

// Reset resets the database to a clean state
func Reset() error {
	fmt.Println("Resetting database...")

	// Create backup before reset
	if err := Backup(); err != nil {
		fmt.Printf("Warning: Failed to create backup before reset: %v\n", err)
	}

	// Remove existing database
	if err := Clean(); err != nil {
		return fmt.Errorf("failed to clean existing database: %v", err)
	}

	// Reinitialize database
	if err := InitDB(); err != nil {
		return fmt.Errorf("failed to reinitialize database: %v", err)
	}

	fmt.Println("Database reset successfully")
	return nil
}

// Status checks the database status
func Status() error {
	// Check if database file exists
	info, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		fmt.Println("Database status: Not initialized")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to check database status: %v", err)
	}

	// Try to open and query the database
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Get table counts
	var stats struct {
		words      int
		groups     int
		activities int
		sessions   int
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&stats.words); err != nil {
		return fmt.Errorf("failed to query words count: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&stats.groups); err != nil {
		return fmt.Errorf("failed to query groups count: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&stats.activities); err != nil {
		return fmt.Errorf("failed to query activities count: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&stats.sessions); err != nil {
		return fmt.Errorf("failed to query sessions count: %v", err)
	}

	fmt.Printf("Database status: Initialized\n")
	fmt.Printf("File size: %d bytes\n", info.Size())
	fmt.Printf("Last modified: %s\n", info.ModTime().Format(time.RFC3339))
	fmt.Printf("Content summary:\n")
	fmt.Printf("- Words: %d\n", stats.words)
	fmt.Printf("- Groups: %d\n", stats.groups)
	fmt.Printf("- Study Activities: %d\n", stats.activities)
	fmt.Printf("- Study Sessions: %d\n", stats.sessions)

	return nil
}

// Run starts the server
func Run() error {
	fmt.Println("Starting server...")
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes generated files
func Clean() error {
	fmt.Println("Cleaning...")
	if err := os.Remove(dbFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Default is the default target
var Default = Run

// Helper functions

func ensureDBDirectory() error {
	dir := filepath.Dir(dbFile)
	return os.MkdirAll(dir, 0755)
}

func ensureBackupDirectory() error {
	return os.MkdirAll(backupDir, 0755)
}

func executeSQLFile(db *sql.DB, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(content))
	return err
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}