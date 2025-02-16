package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"net/http"
)

// ResetHistory resets all study history
func ResetHistory(c *gin.Context) {
	db, err := models.GetDB()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}

	// Delete study sessions and word reviews
	_, err = db.Exec(`
		DELETE FROM study_sessions;
		DELETE FROM review_items;
	`)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to reset study history")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Study history reset successfully",
	})
}

// FullReset performs a complete reset of the database
func FullReset(c *gin.Context) {
	db, err := models.GetDB()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to database")
		return
	}

	// Drop all tables
	_, err = db.Exec(`
		DROP TABLE IF EXISTS review_items;
		DROP TABLE IF EXISTS study_sessions;
		DROP TABLE IF EXISTS study_activities;
		DROP TABLE IF EXISTS group_items;
		DROP TABLE IF EXISTS groups;
		DROP TABLE IF EXISTS words;
	`)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to drop tables")
		return
	}

	// Run migrations to recreate tables
	err = runMigrations(db)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to recreate tables")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database reset successfully",
	})
}

func runMigrations(db *sql.DB) error {
	// Create tables
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			japanese TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS group_items (
			word_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			PRIMARY KEY (word_id, group_id),
			FOREIGN KEY (group_id) REFERENCES groups(id),
			FOREIGN KEY (word_id) REFERENCES words(id)
		);

		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			activity_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (activity_id) REFERENCES study_activities(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);

		CREATE TABLE IF NOT EXISTS study_session_words (
			session_id INTEGER NOT NULL,
			word_id INTEGER NOT NULL,
			PRIMARY KEY (session_id, word_id),
			FOREIGN KEY (session_id) REFERENCES study_sessions(id),
			FOREIGN KEY (word_id) REFERENCES words(id)
		);

		CREATE TABLE IF NOT EXISTS review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			session_id INTEGER NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (session_id) REFERENCES study_sessions(id)
		);
	`)
	return err
}