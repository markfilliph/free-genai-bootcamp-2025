package service

import (
	"fmt"
	"sort"
	"time"

	"lang-portal/internal/models"
)

// DashboardService handles business logic for dashboard statistics
type DashboardService struct{}

// NewDashboardService creates a new DashboardService
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// DashboardStats contains aggregated statistics for the dashboard
type DashboardStats struct {
	TotalWords         int                           `json:"total_words"`
	TotalGroups        int                           `json:"total_groups"`
	TotalSessions      int                           `json:"total_sessions"`
	WordsLearned       int                           `json:"words_learned"`
	AccuracyRate       float64                       `json:"accuracy_rate"`
	StudyStreak        int                           `json:"study_streak"`
	RecentProgress     []*models.StudySession        `json:"recent_progress"`
	TopGroups          []GroupStats                  `json:"top_groups"`
	WeeklyActivity     []DailyActivity              `json:"weekly_activity"`
}

// GroupStats contains statistics for a group
type GroupStats struct {
	GroupID       int64   `json:"group_id"`
	Name         string  `json:"name"`
	TotalWords   int     `json:"total_words"`
	WordsLearned int     `json:"words_learned"`
	AccuracyRate float64 `json:"accuracy_rate"`
	LastStudied  string  `json:"last_studied"`
}

// DailyActivity represents study activity for a single day
type DailyActivity struct {
	Date          string `json:"date"`
	SessionCount  int    `json:"session_count"`
	WordsReviewed int    `json:"words_reviewed"`
	StudyMinutes  int    `json:"study_minutes"`
}

// GetDashboardStats retrieves aggregated statistics for the dashboard
func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	// Get total words count
	totalWords, err := models.GetTotalWordsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to count words: %v", err)
	}
	stats.TotalWords = totalWords

	// Get total groups count
	_, total, err := models.GetGroups(1, 10, "")
	if err != nil {
		return nil, fmt.Errorf("failed to count groups: %v", err)
	}
	stats.TotalGroups = total

	// Get sessions from the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	sessions, err := s.getSessionsSince(thirtyDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %v", err)
	}
	stats.TotalSessions = len(sessions)

	// Calculate words learned and accuracy
	studiedWords, err := models.GetTotalStudiedWordsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get studied words: %v", err)
	}
	stats.WordsLearned = studiedWords

	var totalReviews, correctReviews int
	for _, session := range sessions {
		reviews, err := models.GetSessionReviews(session.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get session reviews: %v", err)
		}
		for _, review := range reviews {
			totalReviews++
			if review.Correct {
				correctReviews++
			}
		}
	}

	if totalReviews > 0 {
		stats.AccuracyRate = float64(correctReviews) / float64(totalReviews)
	}

	// Calculate study streak
	streak, err := s.calculateStudyStreak()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate streak: %v", err)
	}
	stats.StudyStreak = streak

	// Get recent progress (last 5 sessions)
	recent, err := s.getRecentProgress(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent progress: %v", err)
	}
	stats.RecentProgress = recent

	// Get top performing groups
	topGroups, err := s.getTopGroups(5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top groups: %v", err)
	}
	stats.TopGroups = topGroups

	// Get weekly activity
	weeklyActivity, err := s.getWeeklyActivity()
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly activity: %v", err)
	}
	stats.WeeklyActivity = weeklyActivity

	return &stats, nil
}

// GetGroupDashboard retrieves dashboard statistics for a specific group
func (s *DashboardService) GetGroupDashboard(groupID int64) (*GroupStats, error) {
	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	// Get group words
	words, _, err := models.GetGroupWords(groupID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("error getting group words: %v", err)
	}

	// Get group sessions
	sessions, err := models.GetStudySessionsByGroup(groupID, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get group sessions: %v", err)
	}

	// Calculate statistics
	var totalReviews, correctReviews int
	wordMap := make(map[int64]bool)

	for _, session := range sessions {
		reviews, err := models.GetSessionReviews(session.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get session reviews: %v", err)
		}
		for _, review := range reviews {
			totalReviews++
			if review.Correct {
				correctReviews++
			}
			wordMap[review.WordID] = true
		}
	}

	stats := &GroupStats{
		GroupID:     groupID,
		Name:       group.Name,
		TotalWords: len(words),
		WordsLearned: len(wordMap),
	}

	if totalReviews > 0 {
		stats.AccuracyRate = float64(correctReviews) / float64(totalReviews)
	}

	// Get last study time
	if len(sessions) > 0 {
		stats.LastStudied = sessions[0].CreatedAt.Format(time.RFC3339)
	}

	return stats, nil
}

// Helper functions

func (s *DashboardService) calculateStudyStreak() (int, error) {
	streak := 0
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)

	for i := 0; i < 365; i++ { // Check up to a year back
		date := currentDate.AddDate(0, 0, -i)
		sessions, err := s.getSessionsForDate(date)
		if err != nil {
			return 0, fmt.Errorf("failed to check activity: %v", err)
		}
		if len(sessions) == 0 {
			break
		}
		streak++
	}

	return streak, nil
}

func (s *DashboardService) getRecentProgress(limit int) ([]*models.StudySession, error) {
	sessions, err := models.GetStudySessions(0, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent sessions: %v", err)
	}
	return sessions, nil
}

func (s *DashboardService) getTopGroups(limit int) ([]GroupStats, error) {
	// Get all groups
	groups, _, err := models.GetGroups(1, 10, "")
	if err != nil {
		return nil, fmt.Errorf("error getting groups: %v", err)
	}

	// Calculate stats for each group
	var stats []GroupStats
	for _, group := range groups {
		stat, err := s.GetGroupDashboard(group.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get group stats: %v", err)
		}
		stats = append(stats, *stat)
	}

	// Sort by words learned and accuracy
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].WordsLearned == stats[j].WordsLearned {
			return stats[i].AccuracyRate > stats[j].AccuracyRate
		}
		return stats[i].WordsLearned > stats[j].WordsLearned
	})

	if len(stats) > limit {
		stats = stats[:limit]
	}
	return stats, nil
}

func (s *DashboardService) getWeeklyActivity() ([]DailyActivity, error) {
	endDate := time.Now().UTC().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -6) // Last 7 days

	var result []DailyActivity
	currentDate := startDate

	for !currentDate.After(endDate) {
		sessions, err := s.getSessionsForDate(currentDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get sessions: %v", err)
		}

		activity := DailyActivity{
			Date:         currentDate.Format("2006-01-02"),
			SessionCount: len(sessions),
		}

		// Calculate words reviewed and study time
		for _, session := range sessions {
			reviews, err := models.GetSessionReviews(session.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get session reviews: %v", err)
			}
			activity.WordsReviewed += len(reviews)
			// Estimate study time: 30 seconds per word
			activity.StudyMinutes += (len(reviews) * 30) / 60
		}

		result = append(result, activity)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return result, nil
}

func (s *DashboardService) getSessionsSince(since time.Time) ([]*models.StudySession, error) {
	// Get all sessions with pagination
	var allSessions []*models.StudySession
	page := 0
	limit := 100

	for {
		sessions, err := models.GetStudySessions(page*limit, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get sessions: %v", err)
		}
		if len(sessions) == 0 {
			break
		}

		for _, session := range sessions {
			if session.CreatedAt.Before(since) {
				return allSessions, nil
			}
			allSessions = append(allSessions, session)
		}

		page++
	}

	return allSessions, nil
}

func (s *DashboardService) getSessionsForDate(date time.Time) ([]*models.StudySession, error) {
	sessions, err := s.getSessionsSince(date)
	if err != nil {
		return nil, err
	}

	var result []*models.StudySession
	nextDate := date.AddDate(0, 0, 1)

	for _, session := range sessions {
		if session.CreatedAt.After(date) && session.CreatedAt.Before(nextDate) {
			result = append(result, session)
		}
	}

	return result, nil
}
