package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"lang-portal/internal/models"
	"lang-portal/internal/service"
)

// DashboardHandler handles dashboard-related requests
type DashboardHandler struct {
	dashboardService *service.DashboardService
	sessionService  *service.StudySessionService
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(ds *service.DashboardService, ss *service.StudySessionService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: ds,
		sessionService:  ss,
	}
}

// GetDashboardStats returns aggregated statistics for the dashboard
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get dashboard stats: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetLastStudySession returns information about the most recent study session
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	// Get group ID from query parameter
	groupIDStr := c.Query("group_id")
	var groupID int64
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid group ID")
			return
		}
	}

	session, err := h.sessionService.GetLastGroupSession(groupID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get last study session: "+err.Error())
		return
	}

	if session == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "No study sessions found",
		})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudyProgress returns study progress statistics
func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	// Default to last 30 days if not specified
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}

	since := time.Now().AddDate(0, 0, -days)
	progress, err := h.sessionService.GetUserProgress(since)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study progress: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats returns quick overview statistics
func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	// Get total study sessions
	const maxSessions = 100
	sessions, err := h.sessionService.GetStudySessions(0, maxSessions)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get study sessions")
		return
	}

	// Get total active groups
	groups, err := h.dashboardService.GetGroups(0, 0) // Pass 0 for limit to get all groups
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get groups")
		return
	}

	// Calculate success rate from all word reviews
	var totalCorrect, totalReviews int
	for _, session := range sessions {
		stats, err := h.sessionService.GetStudySessionStats(session.ID)
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