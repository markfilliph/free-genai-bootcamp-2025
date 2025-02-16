package models

import (
	"database/sql"
	"time"
)

type StudyActivity struct {
	ID             int64     `json:"id"`
	StudySessionID int64     `json:"study_session_id,omitempty"`
	GroupID        int64     `json:"group_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetStudyActivity returns study activity by ID
func GetStudyActivity(id int64) (*StudyActivity, error) {
	var activity StudyActivity
	err := DB.QueryRow(`
		SELECT id, study_session_id, group_id, created_at
		FROM study_activities
		WHERE id = ?
	`, id).Scan(&activity.ID, &activity.StudySessionID, &activity.GroupID, &activity.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &activity, nil
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(groupID int64) (*StudyActivity, error) {
	var activity StudyActivity
	err := DB.QueryRow(`
		INSERT INTO study_activities (group_id, created_at)
		VALUES (?, CURRENT_TIMESTAMP)
		RETURNING id, study_session_id, group_id, created_at
	`, groupID).Scan(&activity.ID, &activity.StudySessionID, &activity.GroupID, &activity.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

// GetRecentStudyActivities returns recent study activities for a group
func GetRecentStudyActivities(groupID int64, limit int) ([]StudyActivity, error) {
	rows, err := DB.Query(`
		SELECT id, study_session_id, group_id, created_at
		FROM study_activities
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []StudyActivity
	for rows.Next() {
		var activity StudyActivity
		err := rows.Scan(&activity.ID, &activity.StudySessionID, &activity.GroupID, &activity.CreatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// GetStudyActivityStats returns statistics for a study activity
func GetStudyActivityStats(activityID int64) (map[string]int, error) {
	var totalSessions, totalWords int
	var successRate float64

	// Get total sessions
	err := DB.QueryRow(`
		SELECT COUNT(DISTINCT s.id)
		FROM study_sessions s
		WHERE s.study_activity_id = ?
	`, activityID).Scan(&totalSessions)
	if err != nil {
		return nil, err
	}

	// Get total words and success rate
	err = DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT wr.word_id) as total_words,
			COALESCE(AVG(CASE WHEN wr.correct = 1 THEN 100 ELSE 0 END), 0) as success_rate
		FROM word_review_items wr
		JOIN study_sessions s ON wr.study_session_id = s.id
		WHERE s.study_activity_id = ?
	`, activityID).Scan(&totalWords, &successRate)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_sessions": totalSessions,
		"total_words":    totalWords,
		"success_rate":   int(successRate),
	}, nil
}