package database

import "testing"

func TestDatabaseConnection(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := New(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection with a simple query
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	t.Logf("Successfully connected to MySQL version: %s", version)
}
