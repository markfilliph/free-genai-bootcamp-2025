package models

import (
	"database/sql"
	"fmt"
	"time"
)

// StudyActivity represents a study activity for a group of words
type StudyActivity struct {
	ID           *int64     `json:"id"`
	GroupID      *int64     `json:"group_id"`
	ActivityType *string    `json:"activity_type"`
	CreatedAt    *time.Time `json:"created_at"`
	LastUsedAt   *time.Time `json:"last_used_at"`
}

// StudyActivityWithStats includes study activity data with additional statistics
type StudyActivityWithStats struct {
	StudyActivity
	TotalSessions     *int     `json:"total_sessions"`
	TotalWords        *int     `json:"total_words"`
	AverageAccuracy   *float64 `json:"average_accuracy"`
	LastStudiedAt     *string  `json:"last_studied_at,omitempty"`
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(activity *StudyActivity) error {
	result, err := GetDB().Exec(`
		INSERT INTO study_activities (group_id, activity_type, created_at, last_used_at)
		VALUES (?, ?, ?, ?)`,
		activity.GroupID, activity.ActivityType, activity.CreatedAt, activity.LastUsedAt)
	if err != nil {
		return fmt.Errorf("error creating study activity: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	activity.ID = &id
	return nil
}

// GetStudyActivity retrieves a study activity by ID
func GetStudyActivity(activityID int64) (*StudyActivity, error) {
	var activity StudyActivity
	err := GetDB().QueryRow(`
		SELECT id, group_id, activity_type, created_at, last_used_at
		FROM study_activities 
		WHERE id = ?`, activityID).Scan(
		&activity.ID, &activity.GroupID, &activity.ActivityType,
		&activity.CreatedAt, &activity.LastUsedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("study activity not found: %d", activityID)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying study activity: %v", err)
	}

	return &activity, nil
}

// GetStudyActivityWithStats retrieves a study activity by ID with stats
func GetStudyActivityWithStats(activityID int64) (*StudyActivityWithStats, error) {
	db := GetDB()

	// First get the basic activity info
	var activity StudyActivityWithStats
	err := db.QueryRow(`
		SELECT id, group_id, activity_type, created_at, last_used_at 
		FROM study_activities 
		WHERE id = ?`, activityID).Scan(&activity.ID, &activity.GroupID, &activity.ActivityType, &activity.CreatedAt, &activity.LastUsedAt)
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

// UpdateStudyActivity updates an existing study activity
func UpdateStudyActivity(activity *StudyActivity) error {
	result, err := GetDB().Exec(`
		UPDATE study_activities
		SET activity_type = ?, last_used_at = ?
		WHERE id = ?`,
		activity.ActivityType, activity.LastUsedAt, activity.ID)
	if err != nil {
		return fmt.Errorf("error updating study activity: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("study activity not found: %d", *activity.ID)
	}

	return nil
}

// DeleteStudyActivity deletes a study activity and its sessions
func DeleteStudyActivity(id int64) error {
	// Start transaction
	tx, err := GetDB().Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Delete word reviews
	_, err = tx.Exec(`
		DELETE wr FROM word_reviews wr
		JOIN study_sessions ss ON wr.study_session_id = ss.id
		WHERE ss.study_activity_id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting word reviews: %v", err)
	}

	// Delete study sessions
	_, err = tx.Exec(`DELETE FROM study_sessions WHERE study_activity_id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting study sessions: %v", err)
	}

	// Delete study activity
	result, err := tx.Exec(`DELETE FROM study_activities WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting study activity: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("study activity not found: %d", id)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetStudyActivities retrieves a paginated list of study activities for a group
func GetStudyActivities(groupID int64, page, pageSize int) ([]*StudyActivity, int, error) {
	offset := (page - 1) * pageSize

	// Get activities
	rows, err := GetDB().Query(`
		SELECT SQL_CALC_FOUND_ROWS 
			id, group_id, activity_type, created_at, last_used_at
		FROM study_activities
		WHERE group_id = ?
		ORDER BY last_used_at DESC
		LIMIT ? OFFSET ?`,
		groupID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying study activities: %v", err)
	}
	defer rows.Close()

	var activities []*StudyActivity
	for rows.Next() {
		var a StudyActivity
		if err := rows.Scan(
			&a.ID, &a.GroupID, &a.ActivityType,
			&a.CreatedAt, &a.LastUsedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning study activity: %v", err)
		}
		activities = append(activities, &a)
	}

	// Get total count
	var total int
	err = GetDB().QueryRow("SELECT FOUND_ROWS()").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	return activities, total, nil
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