package service

import (
	"database/sql"
	"lang-portal/internal/models"
	"time"
)

// DashboardService handles business logic for dashboard operations
type DashboardService struct {
	db *sql.DB
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(db *sql.DB) *DashboardService {
	return &DashboardService{
		db: db,
	}
}

// GetLastStudySession returns the most recent study session with additional details
func (s *DashboardService) GetLastStudySession() (*models.StudySessionResponse, error) {
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

	return &models.StudySessionResponse{
		ID:                session.ID,
		GroupID:          session.GroupID,
		GroupName:        group.Name,
		CreatedAt:        session.CreatedAt,
		StudyActivityID:  session.StudyActivityID,
		TotalWords:       stats["total"],
		CorrectWords:     stats["correct"],
	}, nil
}

// GetStudyProgress returns overall study progress statistics
func (s *DashboardService) GetStudyProgress() (*models.StudyProgressResponse, error) {
	// Get total words count
	totalWords, err := models.GetTotalWordsCount()
	if err != nil {
		return nil, err
	}

	// Get total studied words count
	totalStudied, err := models.GetTotalStudiedWordsCount()
	if err != nil {
		return nil, err
	}

	// Calculate remaining and completion percentage
	remaining := totalWords - totalStudied
	var completionPercentage float64
	if totalWords > 0 {
		completionPercentage = (float64(totalStudied) / float64(totalWords)) * 100
	}

	return &models.StudyProgressResponse{
		TotalWords:           totalWords,
		TotalWordsStudied:    totalStudied,
		RemainingWords:       remaining,
		CompletionPercentage: completionPercentage,
	}, nil
}

// GetQuickStats returns quick overview statistics
func (s *DashboardService) GetQuickStats() (*models.QuickStatsResponse, error) {
	// Get success rate from recent reviews
	successRate, err := models.GetRecentReviewsSuccessRate()
	if err != nil {
		return nil, err
	}

	// Get total study sessions
	totalSessions, err := models.GetTotalStudySessionsCount()
	if err != nil {
		return nil, err
	}

	// Get total active groups (groups with study activity in last 30 days)
	activeGroups, err := models.GetActiveGroupsCount(30)
	if err != nil {
		return nil, err
	}

	// Calculate study streak
	streakDays, err := s.calculateStudyStreak()
	if err != nil {
		return nil, err
	}

	return &models.QuickStatsResponse{
		SuccessRate:        successRate,
		TotalStudySessions: totalSessions,
		TotalActiveGroups:  activeGroups,
		StudyStreakDays:    streakDays,
	}, nil
}

// Helper function to calculate study streak
func (s *DashboardService) calculateStudyStreak() (int, error) {
	// Get study sessions ordered by date
	sessions, err := models.GetStudySessionsByDate()
	if err != nil {
		return 0, err
	}

	if len(sessions) == 0 {
		return 0, nil
	}

	// Start from today and go backwards
	today := time.Now().UTC().Truncate(24 * time.Hour)
	streakDays := 0
	lastStudyDate := today

	// Check if studied today
	hasStudiedToday := false
	for _, session := range sessions {
		sessionDate := session.CreatedAt.UTC().Truncate(24 * time.Hour)
		if sessionDate.Equal(today) {
			hasStudiedToday = true
			break
		}
	}

	// If not studied today, start counting from yesterday
	if !hasStudiedToday {
		lastStudyDate = today.AddDate(0, 0, -1)
	}

	// Count consecutive days
	for _, session := range sessions {
		sessionDate := session.CreatedAt.UTC().Truncate(24 * time.Hour)
		if sessionDate.Equal(lastStudyDate) {
			streakDays++
			lastStudyDate = sessionDate.AddDate(0, 0, -1)
		} else if sessionDate.Before(lastStudyDate) {
			break
		}
	}

	return streakDays, nil
}