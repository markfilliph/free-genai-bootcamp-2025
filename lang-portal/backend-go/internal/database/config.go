package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB initializes and returns a database connection
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"root",           // username
		"Mfs1985+",      // password
		"localhost",      // host
		"3306",          // port
		"lang_portal",   // database name
	)

	var dbErr error
	once.Do(func() {
		db, dbErr = sql.Open("mysql", dsn)
		if dbErr != nil {
			return
		}

		dbErr = db.Ping()
		if dbErr != nil {
			return
		}

		dbErr = createTables()
	})

	if dbErr != nil {
		return fmt.Errorf("failed to initialize database: %v", dbErr)
	}

	return nil
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
		return db.Close()
	}
	return nil
}

// createTables creates all necessary database tables if they don't exist
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS word_groups (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS words (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			word VARCHAR(255) NOT NULL,
			meaning TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS study_activities (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS study_sessions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			group_id BIGINT NOT NULL,
			study_activity_id BIGINT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES word_groups(id) ON DELETE CASCADE,
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS word_reviews (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			study_session_id BIGINT NOT NULL,
			word_id BIGINT NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}
