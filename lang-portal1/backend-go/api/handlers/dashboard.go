package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/models"
)

func GetLastStudySession(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := models.GetLastStudySession(db)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusOK, nil)
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get last study session")
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func GetStudyProgress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		progress, err := models.GetStudyProgress(db)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get study progress")
			return
		}

		c.JSON(http.StatusOK, progress)
	}
}

func GetQuickStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := models.GetQuickStats(db)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get quick stats")
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}
