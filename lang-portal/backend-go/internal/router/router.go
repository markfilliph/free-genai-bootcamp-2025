package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"lang-portal/internal/handlers"
	"time"
)

// SetupRouter creates and configures a new Gin router
func SetupRouter() *gin.Engine {
	// Create default gin router with Logger and Recovery middleware
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:          12 * time.Hour,
	}))

	// Setup routes
	setupRoutes(r)

	return r
}

// setupRoutes configures all the routes for our application
func setupRoutes(r *gin.Engine) {
	// API group
	api := r.Group("/api")

	// Dashboard routes
	api.GET("/dashboard/last_study_session", handlers.GetLastStudySession)
	api.GET("/dashboard/study_progress", handlers.GetStudyProgress)
	api.GET("/dashboard/quick-stats", handlers.GetQuickStats)

	// Study activities routes
	api.GET("/study_activities/:id", handlers.GetStudyActivity)
	api.GET("/study_activities/:id/study_sessions", handlers.GetStudyActivitySessions)
	api.POST("/study_activities", handlers.CreateStudyActivity)

	// Words routes
	api.GET("/words", handlers.GetWords)
	api.GET("/words/:id", handlers.GetWord)

	// Groups routes
	api.GET("/groups", handlers.GetGroups)
	api.GET("/groups/:id", handlers.GetGroup)
	api.GET("/groups/:id/words", handlers.GetGroupWords)
	api.GET("/groups/:id/study_sessions", handlers.GetGroupStudySessions)

	// Study sessions routes
	api.GET("/study_sessions", handlers.GetStudySessions)
	api.GET("/study_sessions/:id", handlers.GetStudySession)
	api.GET("/study_sessions/:id/words", handlers.GetStudySessionWords)
	api.POST("/study_sessions/:id/words/:word_id/review", handlers.ReviewWord)

	// Reset routes
	api.POST("/reset_history", handlers.ResetHistory)
	api.POST("/full_reset", handlers.FullReset)
}
