package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/models"
)

func GetWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, perPage := getPaginationParams(c)

		words, total, err := models.GetWords(db, page, perPage)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to get words")
			return
		}

		c.JSON(http.StatusOK, newPaginatedResponse(words, page, total, perPage))
	}
}

func GetWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid word ID")
			return
		}

		word, err := models.GetWord(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(c, http.StatusNotFound, "Word not found")
				return
			}
			respondWithError(c, http.StatusInternalServerError, "Failed to get word")
			return
		}

		c.JSON(http.StatusOK, word)
	}
}
