//go:build mage

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type DB mg.Namespace

// Migrate runs all database migrations
func (DB) Migrate() error {
	fmt.Println("Running migrations...")
	
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Open database connection
	dbPath := filepath.Join(wd, "words.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Read migration files
	migrationsDir := filepath.Join(wd, "db", "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	// Execute each migration file in order
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		fmt.Printf("Executing migration: %s\n", file.Name())
		
		// Read migration file
		migrationPath := filepath.Join(migrationsDir, file.Name())
		migration, err := os.ReadFile(migrationPath)
		if err != nil {
			return err
		}

		// Execute migration
		_, err = db.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", file.Name(), err)
		}
	}

	fmt.Println("Migrations completed successfully")
	return nil
}

// Seed imports sample data into the database
func (DB) Seed() error {
	fmt.Println("Seeding database...")
	// TODO: Implement seeding from JSON files
	return nil
}

// Reset removes the database file
func (DB) Reset() error {
	fmt.Println("Resetting database...")
	
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(wd, "words.db")
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Println("Database reset successfully")
	return nil
}

// Build compiles the application
func Build() error {
	fmt.Println("Building application...")
	return sh.Run("go", "build", "-o", "app")
}

// Run starts the application
func Run() error {
	mg.Deps(Build)
	fmt.Println("Starting application...")
	return sh.Run("./app")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	return os.Remove("app")
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...")
}
