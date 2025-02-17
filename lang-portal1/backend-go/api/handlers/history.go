package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to start transaction")
			return
		}
		defer tx.Rollback()

		// Delete all word review items
		_, err = tx.Exec("DELETE FROM word_review_items")
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete word reviews")
			return
		}

		// Delete all study sessions
		_, err = tx.Exec("DELETE FROM study_sessions")
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete study sessions")
			return
		}

		// Delete all study activities
		_, err = tx.Exec("DELETE FROM study_activities")
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to delete study activities")
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Study history has been reset successfully"})
	}
}
