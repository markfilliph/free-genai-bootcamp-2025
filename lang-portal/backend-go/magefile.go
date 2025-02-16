//go:build mage
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/sh"
)

const dbName = "words.db"

// Init initializes the project
func Init() error {
	if err := InitDB(); err != nil {
		return err
	}
	if err := Migrate(); err != nil {
		return err
	}
	return Seed()
}

// InitDB creates a new SQLite database
func InitDB() error {
	// Remove existing database if it exists
	if _, err := os.Stat(dbName); err == nil {
		if err := os.Remove(dbName); err != nil {
			return fmt.Errorf("failed to remove existing database: %v", err)
		}
	}

	// Create new database file
	file, err := os.Create(dbName)
	if err != nil {
		return fmt.Errorf("failed to create database file: %v", err)
	}
	file.Close()

	return nil
}

// Migrate runs database migrations
func Migrate() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	migrations, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find migration files: %v", err)
	}

	for _, migration := range migrations {
		content, err := os.ReadFile(migration)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %v", migration, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", migration, err)
		}
		fmt.Printf("Executed migration: %s\n", migration)
	}

	return nil
}

// Seed runs database seeding
func Seed() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	seeds, err := filepath.Glob("db/seeds/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find seed files: %v", err)
	}

	for _, seed := range seeds {
		content, err := os.ReadFile(seed)
		if err != nil {
			return fmt.Errorf("failed to read seed %s: %v", seed, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute seed %s: %v", seed, err)
		}
		fmt.Printf("Executed seed: %s\n", seed)
	}

	return nil
}

// Build builds the project
func Build() error {
	return sh.Run("go", "build", "-o", "server", "./cmd/server")
}

// Run runs the server
func Run() error {
	return sh.Run("go", "run", "./cmd/server")
}

// Test runs the tests
func Test() error {
	return sh.Run("go", "test", "./...")
}

// Clean cleans build artifacts
func Clean() error {
	return os.Remove("server")
}
