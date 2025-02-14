package main

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"log"
	"net/http"
)

func main() {
	// Initialize database
	if err := models.InitDB("words.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer models.CloseDB()

	// Initialize Gin router
	r := gin.Default()

	// Initialize routes
	initializeRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeRoutes(r *gin.Engine) {
	// API group
	api := r.Group("/api")
	
	// Dashboard routes
	api.GET("/dashboard/last_study_session", getLastStudySession)
	api.GET("/dashboard/study_progress", getStudyProgress)
	api.GET("/dashboard/quick-stats", getQuickStats)

	// Study activities routes
	api.GET("/study_activities/:id", getStudyActivity)
	api.GET("/study_activities/:id/study_sessions", getStudyActivitySessions)
	api.POST("/study_activities", createStudyActivity)

	// Words routes
	api.GET("/words", getWords)
	api.GET("/words/:id", getWord)

	// Groups routes
	api.GET("/groups", getGroups)
	api.GET("/groups/:id", getGroup)
	api.GET("/groups/:id/words", getGroupWords)
	api.GET("/groups/:id/study_sessions", getGroupStudySessions)

	// Study sessions routes
	api.GET("/study_sessions", getStudySessions)
	api.GET("/study_sessions/:id", getStudySession)
	api.GET("/study_sessions/:id/words", getStudySessionWords)
	api.POST("/study_sessions/:id/words/:word_id/review", reviewWord)
}

// Basic handler implementations
func getLastStudySession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudyProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getQuickStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudyActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudyActivitySessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func createStudyActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getWords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getWord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getGroups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getGroupWords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getGroupStudySessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudySessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudySession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func getStudySessionWords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}

func reviewWord(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"})
}