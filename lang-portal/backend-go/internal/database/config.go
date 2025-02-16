package database

import (
	"database/sql"
	"fmt"
	"lang-portal/internal/models"
	"os"
	"sync"
)

var (
	once sync.Once
)

// InitDB initializes and returns a database connection
func InitDB() error {
	var dbErr error
	once.Do(func() {
		dbErr = models.InitDB("")
	})

	if dbErr != nil {
		return fmt.Errorf("failed to initialize database: %v", dbErr)
	}

	return nil
}

// GetDB returns the singleton database connection
func GetDB() (*sql.DB, error) {
	return models.GetDB()
}

// CloseDB closes the database connection
func CloseDB() error {
	return models.CloseDB()
}

// createTables creates all necessary database tables if they don't exist
func createTables() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	// Read and execute the schema file
	schema, err := os.ReadFile("db/migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}
