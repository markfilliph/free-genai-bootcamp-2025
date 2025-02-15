package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lang-portal/internal/handlers"
)

// SetupRouter initializes the Gin router with all routes and middleware
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	setupRoutes(r)

	return r
}

// setupRoutes configures all API routes
func setupRoutes(r *gin.Engine) {
	// Dashboard routes
	r.GET("/api/dashboard/last-study-session", handlers.GetLastStudySession)
	r.GET("/api/dashboard/study-progress", handlers.GetStudyProgress)
	r.GET("/api/dashboard/quick-stats", handlers.GetQuickStats)

	// Study activity routes
	r.GET("/api/study-activities/:id", handlers.GetStudyActivity)
	r.GET("/api/study-activities/:id/sessions", handlers.GetStudyActivitySessions)
	r.POST("/api/study-activities", handlers.CreateStudyActivity)

	// Study session routes
	r.GET("/api/study-sessions", handlers.GetStudySessions)
	r.GET("/api/study-sessions/:id", handlers.GetStudySession)
	r.GET("/api/study-sessions/:id/words", handlers.GetStudySessionWords)
	r.POST("/api/study-sessions/:id/words/:word_id/review", handlers.ReviewWord)

	// Word routes
	r.GET("/api/words", handlers.GetWords)
	r.GET("/api/words/:id", handlers.GetWord)

	// Group routes
	r.GET("/api/groups", handlers.GetGroups)
	r.GET("/api/groups/:id", handlers.GetGroup)
	r.GET("/api/groups/:id/words", handlers.GetGroupWords)
	r.GET("/api/groups/:id/study-sessions", handlers.GetGroupStudySessions)

	// Reset routes
	r.POST("/api/reset/history", handlers.ResetHistory)
	r.POST("/api/reset/full", handlers.FullReset)
}
