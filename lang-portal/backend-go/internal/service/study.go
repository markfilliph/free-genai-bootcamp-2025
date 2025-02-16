package service

import (
	"errors"
	"lang-portal/internal/models"
)

// StudyService handles business logic for study operations
type StudyService struct{}

// NewStudyService creates a new study service
func NewStudyService() *StudyService {
	return &StudyService{}
}

// GetStudyActivities returns all study activities with additional statistics
func (s *StudyService) GetStudyActivities() ([]map[string]interface{}, error) {
	activities, err := models.GetStudyActivities()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, activity := range activities {
		stats, err := models.GetStudyActivityStats(activity.ID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":                activity.ID,
			"name":              "Vocabulary Quiz", // TODO: Add activity types
			"total_sessions":    stats["total_sessions"],
			"total_words":       stats["total_words"],
			"success_rate":      stats["success_rate"],
			"created_at":        activity.CreatedAt,
		})
	}

	return result, nil
}

// GetStudyActivity returns details of a specific study activity
func (s *StudyService) GetStudyActivity(id int64) (map[string]interface{}, error) {
	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errors.New("study activity not found")
	}

	// Get study session associated with this activity
	session, err := models.GetStudySessionByActivityID(id)
	if err != nil {
		return nil, err
	}

	// Get activity stats
	stats, err := models.GetStudyActivityStats(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":               activity.ID,
		"name":             "Vocabulary Quiz", // TODO: Add activity types
		"total_sessions":   stats["total_sessions"],
		"total_words":      stats["total_words"],
		"success_rate":     stats["success_rate"],
		"created_at":       activity.CreatedAt,
		"session":          session,
	}, nil
}

// GetStudyActivitySessions returns study sessions for an activity with detailed statistics
func (s *StudyService) GetStudyActivitySessions(activityID int64) ([]map[string]interface{}, error) {
	// Get all sessions with a large limit since we'll filter by activity ID
	sessions, err := models.GetStudySessions(0, 1000)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, session := range sessions {
		// Skip if activity ID doesn't match or is nil
		if session.StudyActivityID == nil || *session.StudyActivityID != activityID {
			continue
		}

		group, err := models.GetGroup(session.GroupID)
		if err != nil {
			continue
		}

		stats, err := models.GetStudySessionStats(session.ID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":               session.ID,
			"group_name":       group.Name,
			"created_at":       session.CreatedAt,
			"total_words":      stats["total"],
			"correct_words":    stats["correct"],
		})
	}

	return result, nil
}

// CreateStudySession creates a new study session and returns its details
func (s *StudyService) CreateStudySession(groupID, activityID int64) (map[string]interface{}, error) {
	// Create the study session first
	session, err := models.CreateStudySession(groupID)
	if err != nil {
		return nil, err
	}

	// Link the session to the activity if provided
	if activityID != 0 {
		err = models.LinkStudyActivityToSession(activityID, session.ID)
		if err != nil {
			return nil, err
		}
	}

	group, err := models.GetGroup(session.GroupID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"group_name":        group.Name,
		"study_activity_id": session.StudyActivityID,
		"created_at":        session.CreatedAt,
	}, nil
}

// ReviewWord records a word review in a study session
func (s *StudyService) ReviewWord(sessionID, wordID int64, correct bool) (map[string]interface{}, error) {
	review, err := models.CreateWordReview(wordID, sessionID, correct)
	if err != nil {
		return nil, err
	}

	word, err := models.GetWord(wordID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":               review.ID,
		"word_id":          review.WordID,
		"word":             word.Japanese,
		"study_session_id": review.StudySessionID,
		"correct":          review.Correct,
		"created_at":       review.CreatedAt,
	}, nil
}