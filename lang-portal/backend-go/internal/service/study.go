package service

import (
	"database/sql"
	"errors"
	"lang-portal/internal/models"
)

// ActivityType represents the type of study activity
type ActivityType string

const (
	ActivityTypeVocabQuiz ActivityType = "vocabulary_quiz"
	ActivityTypeKanjiDrill ActivityType = "kanji_drill"
)

// StudyService handles business logic for study operations
type StudyService struct {
	db *sql.DB
}

// NewStudyService creates a new study service
func NewStudyService(db *sql.DB) *StudyService {
	return &StudyService{
		db: db,
	}
}

// GetStudyActivities returns all study activities with additional statistics
func (s *StudyService) GetStudyActivities() ([]models.StudyActivityResponse, error) {
	activities, err := models.GetStudyActivities()
	if err != nil {
		return nil, err
	}

	var responses []models.StudyActivityResponse
	for _, activity := range activities {
		stats, err := models.GetStudyActivityStats(activity.ID)
		if err != nil {
			// Use default stats if error
			stats = map[string]int{
				"total_sessions": 0,
				"total_words":    0,
				"success_rate":   0,
			}
		}

		responses = append(responses, models.StudyActivityResponse{
			ID:            activity.ID,
			Name:          s.getActivityTypeName(activity.Type),
			TotalSessions: stats["total_sessions"],
			TotalWords:    stats["total_words"],
			SuccessRate:   float64(stats["success_rate"]),
			CreatedAt:     activity.CreatedAt,
		})
	}

	return responses, nil
}

// GetStudyActivity returns details of a specific study activity
func (s *StudyService) GetStudyActivity(id int64) (*models.StudyActivityResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid activity ID")
	}

	activity, err := models.GetStudyActivity(id)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errors.New("study activity not found")
	}

	// Get study session associated with this activity
	session, err := models.GetStudySessionByActivityID(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Get activity stats
	stats, err := models.GetStudyActivityStats(id)
	if err != nil {
		// Use default stats if error
		stats = map[string]int{
			"total_sessions": 0,
			"total_words":    0,
			"success_rate":   0,
		}
	}

	return &models.StudyActivityResponse{
		ID:            activity.ID,
		Name:          s.getActivityTypeName(activity.Type),
		TotalSessions: stats["total_sessions"],
		TotalWords:    stats["total_words"],
		SuccessRate:   float64(stats["success_rate"]),
		CreatedAt:     activity.CreatedAt,
		Session:       session,
	}, nil
}

// GetStudyActivitySessions returns study sessions for an activity
func (s *StudyService) GetStudyActivitySessions(activityID int64) ([]models.StudySessionResponse, error) {
	if activityID <= 0 {
		return nil, errors.New("invalid activity ID")
	}

	// First verify activity exists
	activity, err := models.GetStudyActivity(activityID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errors.New("study activity not found")
	}

	// Get sessions directly by activity ID
	sessions, err := models.GetStudySessionsByActivityID(activityID)
	if err != nil {
		return nil, err
	}

	var responses []models.StudySessionResponse
	for _, session := range sessions {
		group, err := models.GetGroup(session.GroupID)
		if err != nil {
			continue
		}

		stats, err := models.GetStudySessionStats(session.ID)
		if err != nil {
			continue
		}

		responses = append(responses, models.StudySessionResponse{
			ID:              session.ID,
			GroupID:         session.GroupID,
			GroupName:       group.Name,
			StudyActivityID: &activityID,
			TotalWords:      stats["total"],
			CorrectWords:    stats["correct"],
			CreatedAt:       session.CreatedAt,
		})
	}

	return responses, nil
}

// CreateStudyActivity creates a new study activity
func (s *StudyService) CreateStudyActivity(groupID int64, activityType ActivityType) (*models.StudyActivity, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	// Verify group exists
	group, err := models.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	// Verify activity type
	if !s.isValidActivityType(activityType) {
		return nil, errors.New("invalid activity type")
	}

	return models.CreateStudyActivity(groupID, string(activityType))
}

// ReviewWord records a word review in a study session
func (s *StudyService) ReviewWord(sessionID, wordID int64, correct bool) error {
	if sessionID <= 0 {
		return errors.New("invalid session ID")
	}
	if wordID <= 0 {
		return errors.New("invalid word ID")
	}

	// Verify session exists
	session, err := models.GetStudySession(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("study session not found")
	}

	// Verify word exists
	word, err := models.GetWord(wordID)
	if err != nil {
		return err
	}
	if word == nil {
		return errors.New("word not found")
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create review
	if err := models.CreateWordReview(tx, sessionID, wordID, correct); err != nil {
		return err
	}

	return tx.Commit()
}

// Helper functions

func (s *StudyService) getActivityTypeName(activityType string) string {
	switch ActivityType(activityType) {
	case ActivityTypeVocabQuiz:
		return "Vocabulary Quiz"
	case ActivityTypeKanjiDrill:
		return "Kanji Drill"
	default:
		return "Unknown Activity"
	}
}

func (s *StudyService) isValidActivityType(activityType ActivityType) bool {
	switch activityType {
	case ActivityTypeVocabQuiz, ActivityTypeKanjiDrill:
		return true
	default:
		return false
	}
}