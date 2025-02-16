package handlers

import (
	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"net/http"
	"time"
)

// GetLastStudySession returns information about the most recent study session
func GetLastStudySession(c *gin.Context) {
	session, err := models.GetLastStudySession()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get last study session")
		return
	}

	group, err := models.GetGroup(session.GroupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get group info")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"created_at":        session.CreatedAt,
		"study_activity_id": session.StudyActivityID,
		"group_name":        group.Name,
	})
}

// GetStudyProgress returns study progress statistics
func GetStudyProgress(c *gin.Context) {
	stats, err := models.GetStudyProgress()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study progress")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_words_studied":   stats["total_words_studied"],
		"total_available_words": stats["total_words"],
	})
}

// GetQuickStats returns quick overview statistics
func GetQuickStats(c *gin.Context) {
	// Get total study sessions
	const maxSessions = 100
	sessions, err := models.GetStudySessions(0, maxSessions)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Get total active groups
	groups, err := models.GetGroups()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	// Calculate success rate from all word reviews
	var totalCorrect, totalReviews int
	for _, session := range sessions {
		stats, err := models.GetStudySessionStats(session.ID)
		if err != nil {
			continue
		}
		totalCorrect += stats["correct"]
		totalReviews += stats["total"]
	}

	var successRate float64
	if totalReviews > 0 {
		successRate = float64(totalCorrect) / float64(totalReviews) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"success_rate":         successRate,
		"total_study_sessions": len(sessions),
		"total_active_groups":  len(groups),
		"study_streak_days":    calculateStudyStreak(sessions),
	})
}

// calculateStudyStreak calculates the current study streak in days
func calculateStudyStreak(sessions []*models.StudySession) int {
	if len(sessions) == 0 {
		return 0
	}

	// Sort sessions by date in descending order (most recent first)
	// Note: We assume sessions are already sorted by created_at DESC from the database

	// Get today's date at midnight for comparison
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	streak := 0
	lastDate := today
	
	// Check each day, starting from today
	for i := 0; i <= len(sessions); i++ {
		hasStudyForDay := false
		
		// Look for any sessions on this day
		for _, session := range sessions {
			sessionDate := session.CreatedAt.Truncate(24 * time.Hour)
			if sessionDate.Equal(lastDate) {
				hasStudyForDay = true
				break
			}
		}
		
		if !hasStudyForDay {
			break
		}
		
		streak++
		lastDate = lastDate.AddDate(0, 0, -1) // Go back one day
	}
	
	return streak
}