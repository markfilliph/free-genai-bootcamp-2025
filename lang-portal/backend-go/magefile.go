//go:build mage
package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "words.db"

// Install installs project dependencies
func Install() error {
	fmt.Println("Installing dependencies...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

// InitDB initializes the database with schema and seed data
func InitDB() error {
	fmt.Println("Initializing database...")

	// Create database if it doesn't exist
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// Read and execute schema
	schema, err := os.ReadFile("db/migrations/001_initial_schema.sql")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return err
	}

	// Read and execute seed data
	seeds, err := os.ReadFile("db/seeds/initial_data.sql")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(seeds)); err != nil {
		return err
	}

	fmt.Println("Database initialized successfully")
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