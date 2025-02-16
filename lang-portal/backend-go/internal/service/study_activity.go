package service

import (
	"fmt"
	"time"

	"lang-portal/internal/models"
)

// StudyActivityService handles business logic for study activities
type StudyActivityService struct{}

// NewStudyActivityService creates a new StudyActivityService
func NewStudyActivityService() *StudyActivityService {
	return &StudyActivityService{}
}

// CreateStudyActivity creates a new study activity for a group
func (s *StudyActivityService) CreateStudyActivity(groupID int64, activityType string) (*models.StudyActivity, error) {
	// Validate activity type
	if !isValidActivityType(activityType) {
		return nil, fmt.Errorf("invalid activity type: %s", activityType)
	}

	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	activity := &models.StudyActivity{
		GroupID:      groupID,
		ActivityType: activityType,
		CreatedAt:    time.Now(),
		LastUsedAt:   time.Now(),
		TotalSessions: 0,
	}

	if err := models.CreateStudyActivity(activity); err != nil {
		return nil, fmt.Errorf("failed to create study activity: %v", err)
	}

	return activity, nil
}

// GetStudyActivity retrieves a study activity by ID with optional stats
func (s *StudyActivityService) GetStudyActivity(id int64, withStats bool) (*models.StudyActivity, *models.StudyActivityStats, error) {
	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, nil, fmt.Errorf("study activity not found: %d", id)
	}

	if !withStats {
		return activity, nil, nil
	}

	stats, err := models.GetStudyActivityStats(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get activity stats: %v", err)
	}

	return activity, stats, nil
}

// ListGroupActivities retrieves all study activities for a group
func (s *StudyActivityService) ListGroupActivities(groupID int64, page int) ([]models.StudyActivity, int, error) {
	if page < 1 {
		page = 1
	}

	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get group: %v", err)
	}
	if group == nil {
		return nil, 0, fmt.Errorf("group not found: %d", groupID)
	}

	activities, total, err := models.ListStudyActivities(groupID, page)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list study activities: %v", err)
	}

	return activities, total, nil
}

// UpdateStudyActivity updates an existing study activity
func (s *StudyActivityService) UpdateStudyActivity(id int64, activityType string) (*models.StudyActivity, error) {
	// Validate activity type
	if !isValidActivityType(activityType) {
		return nil, fmt.Errorf("invalid activity type: %s", activityType)
	}

	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found: %d", id)
	}

	activity.ActivityType = activityType
	activity.LastUsedAt = time.Now()

	if err := models.UpdateStudyActivity(activity); err != nil {
		return nil, fmt.Errorf("failed to update study activity: %v", err)
	}

	return activity, nil
}

// DeleteStudyActivity deletes a study activity and its sessions
func (s *StudyActivityService) DeleteStudyActivity(id int64) error {
	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return fmt.Errorf("study activity not found: %d", id)
	}

	if err := models.DeleteStudyActivity(id); err != nil {
		return fmt.Errorf("failed to delete study activity: %v", err)
	}

	return nil
}

// GetActivityProgress retrieves progress statistics for a study activity
func (s *StudyActivityService) GetActivityProgress(id int64) (*models.StudyActivityProgress, error) {
	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found: %d", id)
	}

	progress, err := models.GetStudyActivityProgress(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity progress: %v", err)
	}

	return progress, nil
}

// StartStudySession starts a new study session for an activity
func (s *StudyActivityService) StartStudySession(activityID int64) (*models.StudySession, error) {
	activity, err := models.GetStudyActivity(activityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found: %d", activityID)
	}

	session, err := models.CreateStudySession(activity.GroupID, activityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create study session: %v", err)
	}

	// Update activity stats
	activity.LastUsedAt = time.Now()
	activity.TotalSessions++
	if err := models.UpdateStudyActivity(activity); err != nil {
		return nil, fmt.Errorf("failed to update activity stats: %v", err)
	}

	return session, nil
}

// GetRecentSessions retrieves recent study sessions for an activity
func (s *StudyActivityService) GetRecentSessions(activityID int64, limit int) ([]models.StudySessionWithStats, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	activity, err := models.GetStudyActivity(activityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found: %d", activityID)
	}

	sessions, _, err := models.GetStudySessionsByActivity(activityID, 1) // Get first page
	if err != nil {
		return nil, fmt.Errorf("failed to get recent sessions: %v", err)
	}

	if len(sessions) > limit {
		sessions = sessions[:limit]
	}

	return sessions, nil
}

// Helper functions

func isValidActivityType(activityType string) bool {
	validTypes := map[string]bool{
		"flashcards": true,
		"quiz":       true,
		"writing":    true,
		"listening":  true,
		"speaking":   true,
	}
	return validTypes[activityType]
}
