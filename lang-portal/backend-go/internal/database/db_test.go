package database

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test that we can execute a simple query
	var version string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to query database version: %v", err)
	}
	t.Logf("Successfully connected to SQLite version: %s", version)
}
