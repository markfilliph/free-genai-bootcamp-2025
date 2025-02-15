package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

// Config holds the database configuration
type Config struct {
	// User     string
	// Password string
	// Host     string
	// Port     string
	// DBName   string
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() (*Config, error) {
	// Try to load .env from the project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		// User:     os.Getenv("DB_USER"),
		// Password: os.Getenv("DB_PASSWORD"),
		// Host:     os.Getenv("DB_HOST"),
		// Port:     os.Getenv("DB_PORT"),
		// DBName:   os.Getenv("DB_NAME"),
	}, nil
}

// DSN returns the Data Source Name
func (c *Config) DSN() string {
	// return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
	// 	c.User, c.Password, c.Host, c.Port, c.DBName)
	return ""
}

// GetDBPath returns the absolute path to the SQLite database file
func GetDBPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	return filepath.Join(basepath, "words.db"), nil
}

// InitDB initializes and returns a database connection
func InitDB() (*sql.DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}

	// Open SQLite3 database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
