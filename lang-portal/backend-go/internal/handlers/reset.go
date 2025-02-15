package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"lang-portal/internal/database"
	"net/http"
)

// ResetHistory resets all study history
func ResetHistory(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Delete all study sessions and related data
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete word reviews")
		return
	}

	_, err = tx.Exec("DELETE FROM study_activities")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete study activities")
		return
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete study sessions")
		return
	}

	if err := tx.Commit(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Study history has been reset",
	})
}

// FullReset performs a complete reset of the database
func FullReset(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Drop all tables
	tables := []string{
		"word_review_items",
		"study_activities",
		"study_sessions",
		"words_groups",
		"words",
		"groups",
	}

	for _, table := range tables {
		_, err = tx.Exec("DROP TABLE IF EXISTS " + table)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to drop table "+table)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Re-run migrations
	if err := runMigrations(db); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to run migrations")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database has been reset to initial state",
	})
}

func runMigrations(db *sql.DB) error {
	// Read and execute migration file
	migration := `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		japanese TEXT NOT NULL,
		romaji TEXT NOT NULL,
		english TEXT NOT NULL,
		parts TEXT
	);

	CREATE TABLE IF NOT EXISTS words_groups (
		word_id INTEGER,
		group_id INTEGER,
		PRIMARY KEY (word_id, group_id),
		FOREIGN KEY (word_id) REFERENCES words(id),
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);

	CREATE TABLE IF NOT EXISTS study_activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		study_session_id INTEGER,
		group_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);

	CREATE TABLE IF NOT EXISTS study_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		study_activity_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES groups(id),
		FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
	);

	CREATE TABLE IF NOT EXISTS word_review_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word_id INTEGER NOT NULL,
		study_session_id INTEGER NOT NULL,
		correct BOOLEAN NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (word_id) REFERENCES words(id),
		FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
	);`

	_, err := db.Exec(migration)
	return err
}