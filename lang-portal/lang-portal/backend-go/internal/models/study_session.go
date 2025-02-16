package models

import (
	"database/sql"
	"fmt"
	"time"
)

// StudySession represents a single study session
type StudySession struct {
	ID              *int64     `json:"id"`
	GroupID         *int64     `json:"group_id"`
	StudyActivityID *int64     `json:"study_activity_id"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	CreatedAt       *time.Time `json:"created_at"`
}

// StudySessionWithStats includes study session data with statistics
type StudySessionWithStats struct {
	StudySession
	TotalWords      int     `json:"total_words"`
	CorrectWords    int     `json:"correct_words"`
	IncorrectWords  int     `json:"incorrect_words"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	CompletionRate  float64 `json:"completion_rate"`
	ReviewedWords   []int64 `json:"reviewed_words,omitempty"`
}

// CreateStudySession creates a new study session
func CreateStudySession(session *StudySession) error {
	result, err := GetDB().Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, start_time, end_time, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		session.GroupID, session.StudyActivityID,
		session.StartTime, session.EndTime,
		session.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	session.ID = &id
	return nil
}

// GetStudySession retrieves a study session by ID
func GetStudySession(id int64) (*StudySession, error) {
	var session StudySession
	err := GetDB().QueryRow(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE id = ?`, id).Scan(
		&session.ID, &session.GroupID,
		&session.StudyActivityID, &session.StartTime,
		&session.EndTime, &session.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting study session: %v", err)
	}

	return &session, nil
}

// GetStudySessions retrieves all study sessions with pagination
func GetStudySessions(offset, limit int) ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetStudySessionsByGroup retrieves study sessions for a specific group
func GetStudySessionsByGroup(groupID int64, offset, limit int) ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetStudySessionsByGroupID retrieves all study sessions for a specific group
func GetStudySessionsByGroupID(groupID int64) ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// GetLastStudySession returns the most recent study session
func GetLastStudySession() (*StudySession, error) {
	var s StudySession
	err := GetDB().QueryRow(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at 
		FROM study_sessions 
		ORDER BY created_at DESC 
		LIMIT 1
	`).Scan(&s.ID, &s.GroupID, &s.StudyActivityID, &s.StartTime, &s.EndTime, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpdateStudyActivityID updates the study activity ID for a session
func (s *StudySession) UpdateStudyActivityID(activityID int64) error {
	_, err := GetDB().Exec(`
		UPDATE study_sessions
		SET study_activity_id = ?
		WHERE id = ?
	`, activityID, s.ID)
	if err != nil {
		return err
	}

	s.StudyActivityID = &activityID
	return nil
}

// LinkStudyActivityToSession links a study activity to a study session
func LinkStudyActivityToSession(activityID, studySessionID int64) error {
	tx, err := GetDB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the study activity with the session ID
	_, err = tx.Exec(`
		UPDATE study_activities
		SET study_session_id = ?
		WHERE id = ?
	`, studySessionID, activityID)
	if err != nil {
		return err
	}

	// Update the study session with the activity ID
	_, err = tx.Exec(`
		UPDATE study_sessions
		SET study_activity_id = ?
		WHERE id = ?
	`, activityID, studySessionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetStudySessionStats returns statistics for a study session
func GetStudySessionStats(sessionID int64) (map[string]int, error) {
	var total, correct int
	err := GetDB().QueryRow(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct
		FROM word_review_items
		WHERE study_session_id = ?
	`, sessionID).Scan(&total, &correct)

	if err != nil {
		return nil, err
	}

	var successRate int
	if total > 0 {
		successRate = (correct * 100) / total
	}

	return map[string]int{
		"total_reviews":    total,
		"correct_reviews": correct,
		"success_rate":    successRate,
	}, nil
}

// GetStudySessionByActivityID retrieves a study session by activity ID
func GetStudySessionByActivityID(activityID int64) (*StudySession, error) {
	var session StudySession
	err := GetDB().QueryRow(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		LIMIT 1
	`, activityID).Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetStudySessionsByActivity retrieves study sessions for a specific activity with pagination
func GetStudySessionsByActivity(activityID int64, offset, limit int) ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, activityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// GetTotalStudySessionsByActivity returns the total count of study sessions for an activity
func GetTotalStudySessionsByActivity(activityID int64) (int, error) {
	var count int
	err := GetDB().QueryRow(`
		SELECT COUNT(*)
		FROM study_sessions
		WHERE study_activity_id = ?
	`, activityID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetStudySessionsByActivityID retrieves all study sessions for a specific activity
func GetStudySessionsByActivityID(activityID int64) ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
	`, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetStudySessionsByDate returns all study sessions ordered by date
func GetStudySessionsByDate() ([]*StudySession, error) {
	rows, err := GetDB().Query(`
		SELECT id, group_id, study_activity_id, start_time, end_time, created_at
		FROM study_sessions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*StudySession
	for rows.Next() {
		var session StudySession
		err := rows.Scan(&session.ID, &session.GroupID, &session.StudyActivityID, &session.StartTime, &session.EndTime, &session.CreatedAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetTotalWordsCount returns the total number of words in the system
func GetTotalWordsCount() (int, error) {
	var count int
	err := GetDB().QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalStudiedWordsCount returns the total number of words that have been studied
func GetTotalStudiedWordsCount() (int, error) {
	var count int
	err := GetDB().QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_review_items
	`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalStudySessionsCount returns the total number of study sessions
func GetTotalStudySessionsCount() (int, error) {
	var count int
	err := GetDB().QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetActiveGroupsCount returns the number of groups with study activity in the last N days
func GetActiveGroupsCount(days int) (int, error) {
	var count int
	err := GetDB().QueryRow(`
		SELECT COUNT(DISTINCT g.id)
		FROM word_groups g
		JOIN study_sessions s ON g.id = s.group_id
		WHERE s.created_at >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
	`, days).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}