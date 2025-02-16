package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) error {
	var err error
	
	// Load .env file if it exists
	_ = godotenv.Load()
	
	// If no DSN provided, use environment variables or defaults
	if dataSourceName == "" {
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			getEnvOrDefault("DB_USER", "root"),
			getEnvOrDefault("DB_PASSWORD", ""),
			getEnvOrDefault("DB_HOST", "localhost"),
			getEnvOrDefault("DB_PORT", "3306"),
			getEnvOrDefault("DB_NAME", "lang_portal"),
		)
	}

	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

// GetDB returns the database connection
func GetDB() (*sql.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	return DB, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
