package database

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	if err := InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer CloseDB()

	db, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}

	// Test that we can execute a simple query
	var version string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to query database version: %v", err)
	}
	t.Logf("Successfully connected to SQLite version: %s", version)
}
