package service

import (
	"fmt"
	"time"

	"lang-portal/internal/models"
)

type StudyActivityService struct {
	// Add any dependencies here
}

func NewStudyActivityService() *StudyActivityService {
	return &StudyActivityService{}
}

// CreateStudyActivity creates a new study activity
func (s *StudyActivityService) CreateStudyActivity(groupID int64, activityType string) (*models.StudyActivity, error) {
	now := time.Now()
	gid := groupID
	atype := activityType
	activity := &models.StudyActivity{
		GroupID:      &gid,
		ActivityType: &atype,
		LastUsedAt:   &now,
		CreatedAt:    &now,
	}

	if err := models.CreateStudyActivity(activity); err != nil {
		return nil, fmt.Errorf("error creating study activity: %v", err)
	}

	return activity, nil
}

// GetStudyActivity retrieves a study activity by ID
func (s *StudyActivityService) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, fmt.Errorf("error getting study activity: %v", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found: %d", id)
	}
	return activity, nil
}

// GetStudyActivities retrieves study activities for a group with pagination
func (s *StudyActivityService) GetStudyActivities(groupID int64, offset, limit int) ([]*models.StudyActivity, int, error) {
	activities, total, err := models.GetStudyActivities(groupID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting study activities: %v", err)
	}
	return activities, total, nil
}

// UpdateStudyActivity updates a study activity
func (s *StudyActivityService) UpdateStudyActivity(activity *models.StudyActivity) error {
	now := time.Now()
	activity.LastUsedAt = &now
	if err := models.UpdateStudyActivity(activity); err != nil {
		return fmt.Errorf("error updating study activity: %v", err)
	}
	return nil
}

// DeleteStudyActivity deletes a study activity and its associated data
func (s *StudyActivityService) DeleteStudyActivity(id int64) error {
	if err := models.DeleteStudyActivity(id); err != nil {
		return fmt.Errorf("error deleting study activity: %v", err)
	}
	return nil
}

// GetStudyActivityStats retrieves statistics for a study activity
func (s *StudyActivityService) GetStudyActivityStats(activityID int64) (*models.StudyActivityStats, error) {
	stats, err := models.GetStudyActivityStats(activityID)
	if err != nil {
		return nil, fmt.Errorf("error getting study activity stats: %v", err)
	}
	return stats, nil
}

// GetStudyActivityProgress retrieves progress for a study activity
func (s *StudyActivityService) GetStudyActivityProgress(activityID int64) (*models.StudyActivityProgress, error) {
	progress, err := models.GetStudyActivityProgress(activityID)
	if err != nil {
		return nil, fmt.Errorf("error getting study activity progress: %v", err)
	}
	return progress, nil
}

// GetStudySessionsByActivity retrieves study sessions for an activity with pagination
func (s *StudyActivityService) GetStudySessionsByActivity(activityID int64, offset, limit int) ([]*models.StudySession, error) {
	sessions, err := models.GetStudySessionsByActivity(activityID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting study sessions: %v", err)
	}
	return sessions, nil
}
