package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResetHistory resets all study history
func ResetHistory(c *gin.Context) {
	// TODO: Implement with actual database query
	c.JSON(http.StatusOK, gin.H{
		"message": "Study history has been reset",
	})
}

// FullReset performs a complete reset of the database
func FullReset(c *gin.Context) {
	// TODO: Implement with actual database query
	c.JSON(http.StatusOK, gin.H{
		"message": "Database has been reset to initial state",
	})
}