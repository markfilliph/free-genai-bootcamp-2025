package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"lang-portal/internal/models"
)

func main() {
	// Initialize database connection
	if err := models.InitDB("../../words.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	r := gin.Default()

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/dashboard/study_progress", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/dashboard/quick-stats", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})

		// Study activities routes
		api.GET("/study_activities/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/study_activities/:id/study_sessions", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.POST("/study_activities", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})

		// Words routes
		api.GET("/words", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/words/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})

		// Groups routes
		api.GET("/groups", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/groups/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/groups/:id/words", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/groups/:id/study_sessions", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})

		// Study sessions routes
		api.GET("/study_sessions", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/study_sessions/:id", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.GET("/study_sessions/:id/words", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.POST("/study_sessions/:id/words/:word_id/review", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})

		// Reset routes
		api.POST("/reset_history", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
		api.POST("/full_reset", func(c *gin.Context) {
			c.JSON(501, gin.H{"error": "Not implemented"})
		})
	}

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
