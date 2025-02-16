package models

import (
	"database/sql"
	"fmt"
	"time"
)

// StudyActivity represents a study activity for a group of words
type StudyActivity struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
}

// StudyActivityWithStats includes study activity data with additional statistics
type StudyActivityWithStats struct {
	StudyActivity
	TotalSessions     int     `json:"total_sessions"`
	TotalWords        int     `json:"total_words"`
	AverageAccuracy   float64 `json:"average_accuracy"`
	LastStudiedAt     *string `json:"last_studied_at,omitempty"`
}

// GetStudyActivity retrieves a study activity by ID with stats
func GetStudyActivity(activityID int64) (*StudyActivityWithStats, error) {
	db := GetDB()

	// First get the basic activity info
	var activity StudyActivityWithStats
	err := db.QueryRow(`
		SELECT id, group_id, created_at 
		FROM study_activities 
		WHERE id = ?`, activityID).Scan(&activity.ID, &activity.GroupID, &activity.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying study activity: %v", err)
	}

	// Get statistics
	err = db.QueryRow(`
		SELECT 
			COUNT(DISTINCT ss.id) as total_sessions,
			COUNT(DISTINCT wri.word_id) as total_words,
			COALESCE(AVG(CASE WHEN wri.correct THEN 1.0 ELSE 0.0 END), 0) as average_accuracy,
			MAX(ss.created_at) as last_studied_at
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON sa.id = ss.study_activity_id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE sa.id = ?
		GROUP BY sa.id`, activityID).
		Scan(&activity.TotalSessions, &activity.TotalWords, &activity.AverageAccuracy, &activity.LastStudiedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting study activity stats: %v", err)
	}

	return &activity, nil
}

// CreateStudyActivity creates a new study activity for a group
func CreateStudyActivity(groupID int64) (*StudyActivity, error) {
	db := GetDB()

	// Verify group exists
	exists, err := groupExists(groupID)
	if err != nil {
		return nil, fmt.Errorf("error checking group existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("group not found: %d", groupID)
	}

	// Create study activity
	result, err := db.Exec(`
		INSERT INTO study_activities (group_id, created_at)
		VALUES (?, CURRENT_TIMESTAMP)`,
		groupID)
	if err != nil {
		return nil, fmt.Errorf("error creating study activity: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return &StudyActivity{
		ID:      id,
		GroupID: groupID,
	}, nil
}

// GetStudyActivitySessions retrieves all study sessions for an activity
func GetStudyActivitySessions(activityID int64) ([]StudySession, error) {
	db := GetDB()

	rows, err := db.Query(`
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC`, activityID)
	if err != nil {
		return nil, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var s StudySession
		if err := rows.Scan(&s.ID, &s.GroupID, &s.StudyActivityID, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// groupExists checks if a group exists
func groupExists(groupID int64) (bool, error) {
	db := GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)", groupID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking group existence: %v", err)
	}

	return exists, nil
}