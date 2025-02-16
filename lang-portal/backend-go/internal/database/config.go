package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

var (
	db   *sql.DB
	once sync.Once
)

// Config holds the database configuration
type Config struct {
	DBPath string
}

// NewConfig creates a new database configuration
func NewConfig() (*Config, error) {
	// Try to load .env from the project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		// Ignore error if .env doesn't exist
		fmt.Printf("Warning: .env file not found, using default configuration\n")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		var err error
		dbPath, err = GetDBPath()
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		DBPath: dbPath,
	}, nil
}

// GetDBPath returns the absolute path to the SQLite database file
func GetDBPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	return filepath.Join(basepath, "words.db"), nil
}

// InitDB initializes and returns a database connection
func InitDB() error {
	var initErr error
	once.Do(func() {
		config, err := NewConfig()
		if err != nil {
			initErr = err
			return
		}

		// Open SQLite3 database
		var dbConn *sql.DB
		dbConn, err = sql.Open("sqlite3", config.DBPath)
		if err != nil {
			initErr = err
			return
		}

		// Test the connection
		if err := dbConn.Ping(); err != nil {
			initErr = err
			return
		}

		// Set the global db variable
		db = dbConn

		// Create tables if they don't exist
		if err := createTables(); err != nil {
			initErr = err
			return
		}
	})
	return initErr
}

// GetDB returns the singleton database connection
func GetDB() (*sql.DB, error) {
	if db == nil {
		if err := InitDB(); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
		db = nil
	}
	return nil
}

// createTables creates all necessary database tables if they don't exist
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			japanese TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT,
			group_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id)
		)`,
		`CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			study_activity_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
		)`,
		`CREATE TABLE IF NOT EXISTS word_reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			session_id INTEGER NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (session_id) REFERENCES study_sessions(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
