package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"lang-portal/backend/api/handlers"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// API group
	api := r.Group("/api")

	// Dashboard routes
	api.GET("/dashboard/last_study_session", handlers.GetLastStudySession(db))
	api.GET("/dashboard/study_progress", handlers.GetDailyStudyProgress(db))
	api.GET("/dashboard/quick-stats", handlers.GetQuickStats(db))

	// Study activity routes
	api.GET("/study_activities/:id", handlers.GetStudyActivity(db))
	api.GET("/study_activities/:id/study_sessions", handlers.GetStudyActivitySessions(db))
	api.POST("/study_activities", handlers.CreateStudySession(db))

	// Word routes
	api.GET("/words", handlers.GetWords(db))
	api.POST("/words", handlers.CreateWord(db))
	
	// Single word routes
	wordRoutes := api.Group("/words/:id")
	{
		wordRoutes.GET("", handlers.GetWord(db))
		wordRoutes.PUT("", handlers.UpdateWord(db))
		wordRoutes.DELETE("", handlers.DeleteWord(db))
	}

	// Word-group relationship routes
	wordGroupRoutes := api.Group("/word-groups")
	{
		wordGroupRoutes.POST("/:wordId/:groupId", handlers.AddWordToGroup(db))
		wordGroupRoutes.DELETE("/:wordId/:groupId", handlers.RemoveWordFromGroup(db))
	}

	// Group routes
	api.GET("/groups", handlers.GetGroups(db))
	api.GET("/groups/:id", handlers.GetGroup(db))
	api.GET("/groups/:id/words", handlers.GetGroupWords(db))
	api.POST("/groups", handlers.CreateGroup(db))
	api.PUT("/groups/:id", handlers.UpdateGroup(db))
	api.DELETE("/groups/:id", handlers.DeleteGroup(db))
	api.GET("/groups/:id/study_sessions", handlers.GetGroupStudySessions(db))

	// Study session routes
	api.GET("/study_sessions", handlers.GetStudySessions(db))
	api.GET("/study_sessions/:id", handlers.GetStudySession(db))
	api.POST("/study_sessions/:session_id/words/:word_id/review", handlers.AddWordReview(db))

	// Reset routes
	api.POST("/reset_history", handlers.ResetHistory(db))
	// Full reset is not part of the technical specs, so we'll keep it as not implemented
	api.POST("/full_reset", func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "Not implemented"})
	})
}
