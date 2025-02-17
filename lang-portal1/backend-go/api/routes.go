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
		wordGroupRoutes.POST("/:wordId/:groupId", handlers.AddWordToGroupFromWord(db))
		wordGroupRoutes.DELETE("/:wordId/:groupId", handlers.RemoveWordFromGroupHandler(db))
	}

	// Group routes
	groupRoutes := api.Group("/groups")
	{
		groupRoutes.GET("", handlers.GetGroups(db))
		groupRoutes.POST("", handlers.CreateGroup(db))
		groupRoutes.GET("/:id", handlers.GetGroup(db))
		groupRoutes.PUT("/:id", handlers.UpdateGroup(db))
		groupRoutes.DELETE("/:id", handlers.DeleteGroup(db))
		groupRoutes.GET("/:id/words", handlers.GetGroupWords(db))
		groupRoutes.POST("/:id/words/:wordId", handlers.AddWordToGroup(db))
		groupRoutes.DELETE("/:id/words/:wordId", handlers.RemoveWordFromGroup(db))
		groupRoutes.GET("/:id/study_sessions", handlers.GetGroupStudySessions(db))
	}

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
