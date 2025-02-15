package handlers

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"net/http"
	"strconv"
)

// GetWords returns a paginated list of words
func GetWords(c *gin.Context) {
	pagination := getPaginationParams(c)

	words, err := models.GetWords()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get words")
		return
	}

	// Filter and paginate words
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if end > len(words) {
		end = len(words)
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      words[start:end],
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, len(words)),
	})
}

// GetWord returns details of a specific word
func GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word, err := models.GetWord(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get word")
		return
	}

	c.JSON(http.StatusOK, word)
}