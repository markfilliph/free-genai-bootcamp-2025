package service

import (
	"lang-portal/internal/models"
	"time"
)

// DashboardService handles business logic for dashboard operations
type DashboardService struct{}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetLastStudySession returns the most recent study session with additional details
func (s *DashboardService) GetLastStudySession() (map[string]interface{}, error) {
	session, err := models.GetLastStudySession()
	if err != nil {
		return nil, err
	}

	group, err := models.GetGroup(session.GroupID)
	if err != nil {
		return nil, err
	}

	stats, err := models.GetStudySessionStats(session.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"group_name":        group.Name,
		"created_at":        session.CreatedAt,
		"study_activity_id": session.StudyActivityID,
		"total_words":       stats["total"],
		"correct_words":     stats["correct"],
	}, nil
}

// GetStudyProgress returns overall study progress statistics
func (s *DashboardService) GetStudyProgress() (map[string]interface{}, error) {
	stats, err := models.GetStudyProgress()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_words":          stats["total_words"],
		"total_words_studied":  stats["total_words_studied"],
		"remaining_words":      stats["remaining_words"],
		"completion_percentage": stats["completion_percentage"],
	}, nil
}

// GetQuickStats returns quick overview statistics
func (s *DashboardService) GetQuickStats() (map[string]interface{}, error) {
	// Get all sessions with a large limit
	sessions, err := models.GetStudySessions(0, 1000)
	if err != nil {
		return nil, err
	}

	// Get all groups
	groups, err := models.GetGroups(0, 0) // Pass 0 for limit to get all groups
	if err != nil {
		return nil, err
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

	// Calculate study streak
	streak := calculateStudyStreak(sessions)

	return map[string]interface{}{
		"success_rate":         successRate,
		"total_study_sessions": len(sessions),
		"total_active_groups":  len(groups),
		"study_streak_days":    streak,
	}, nil
}

// calculateStudyStreak calculates the current study streak in days
func calculateStudyStreak(sessions []*models.StudySession) int {
	if len(sessions) == 0 {
		return 0
	}

	// Sort sessions by date (they should already be sorted)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	streak := 0
	lastDate := today

	for _, session := range sessions {
		sessionDate := session.CreatedAt.UTC().Truncate(24 * time.Hour)
		if sessionDate.Equal(lastDate) {
			continue
		}
		if sessionDate.Add(24 * time.Hour).Before(lastDate) {
			break
		}
		streak++
		lastDate = sessionDate
	}

	return streak
}