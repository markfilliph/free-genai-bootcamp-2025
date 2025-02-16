//go:build mage
package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
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

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Get database connection parameters from environment
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASSWORD", "")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "lang_portal")

	fmt.Printf("Using database configuration: user=%s, host=%s, port=%s, dbname=%s\n", dbUser, dbHost, dbPort, dbName)

	// Create database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Execute schema file
	if err := executeSQLFile(db, schemaFile); err != nil {
		return fmt.Errorf("failed to execute schema file: %v", err)
	}

	// Execute seed file if it exists
	if _, err := os.Stat(seedFile); err == nil {
		if err := executeSQLFile(db, seedFile); err != nil {
			return fmt.Errorf("failed to execute seed file: %v", err)
		}
	}

	fmt.Println("Database initialization completed successfully")
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

// Backup creates a backup of the database
func Backup() error {
	if err := ensureBackupDirectory(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", timestamp))

	// Get database connection parameters from environment
	_ = godotenv.Load()
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASSWORD", "")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "lang_portal")

	// Create mysqldump command
	cmd := exec.Command("mysqldump",
		"-h", dbHost,
		"-P", dbPort,
		"-u", dbUser,
		fmt.Sprintf("-p%s", dbPass),
		dbName)

	// Open the output file
	outFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer outFile.Close()

	// Set the output to our file
	cmd.Stdout = outFile

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	fmt.Printf("Backup created successfully: %s\n", backupFile)
	return nil
}

// Clean removes generated files
func Clean() error {
	fmt.Println("Cleaning generated files...")
	return nil
}

// Default is the default target
var Default = Run

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func ensureBackupDirectory() error {
	return os.MkdirAll(backupDir, 0755)
}

func executeSQLFile(db *sql.DB, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Split the SQL file into individual statements
	statements := strings.Split(string(content), ";")

	// Execute each statement
	for _, stmt := range statements {
		// Skip empty statements
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute SQL statement '%s': %v", stmt, err)
		}
	}

	return nil
}