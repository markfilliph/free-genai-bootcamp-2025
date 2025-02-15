package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetWords returns a paginated list of words
func GetWords(c *gin.Context) {
	pagination := getPaginationParams(c)

	// TODO: Implement with actual database query
	c.JSON(http.StatusOK, gin.H{
		"items": []gin.H{
			{
				"id":       1,
				"japanese": "こんにちは",
				"romaji":   "konnichiwa",
				"english":  "hello",
				"parts":    gin.H{"type": "greeting", "formality": "neutral"},
			},
		},
		"pagination": calculatePagination(pagination.Page, pagination.PageSize, 100),
	})
}

// GetWord returns details of a specific word
func GetWord(c *gin.Context) {
	id := c.Param("id")
	
	// TODO: Implement with actual database query
	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"japanese": "こんにちは",
		"romaji":   "konnichiwa",
		"english":  "hello",
		"parts":    gin.H{"type": "greeting", "formality": "neutral"},
	})
}