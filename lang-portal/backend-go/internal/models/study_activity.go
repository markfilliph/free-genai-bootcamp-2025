package models

import "time"

type StudyActivity struct {
	ID             int64     `json:"id"`
	StudySessionID int64     `json:"study_session_id,omitempty"`
	GroupID        int64     `json:"group_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id,omitempty"`
}

// GetLastStudySession returns the most recent study session
func GetLastStudySession() (*StudySession, error) {
	var s StudySession
	err := DB.QueryRow(`
		SELECT id, group_id, created_at, study_activity_id 
		FROM study_sessions 
		ORDER BY created_at DESC 
		LIMIT 1
	`).Scan(&s.ID, &s.GroupID, &s.CreatedAt, &s.StudyActivityID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStudyProgress returns study progress statistics
func GetStudyProgress() (map[string]int, error) {
	var totalWords, studiedWords int
	err := DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalWords)
	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&studiedWords)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_words":          totalWords,
		"total_words_studied":  studiedWords,
		"remaining_words":      totalWords - studiedWords,
		"completion_percentage": (studiedWords * 100) / totalWords,
	}, nil
}

// GetStudyActivity returns study activity by ID
func GetStudyActivity(id int64) (*StudyActivity, error) {
	var a StudyActivity
	err := DB.QueryRow(`
		SELECT id, study_session_id, group_id, created_at 
		FROM study_activities 
		WHERE id = ?
	`, id).Scan(&a.ID, &a.StudySessionID, &a.GroupID, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(groupID int64) (*StudyActivity, error) {
	result, err := DB.Exec(`
		INSERT INTO study_activities (group_id)
		VALUES (?)
	`, groupID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetStudyActivity(id)
}

// GetStudySessions returns all study sessions
func GetStudySessions() ([]StudySession, error) {
	rows, err := DB.Query(`
		SELECT id, group_id, created_at, study_activity_id 
		FROM study_sessions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []StudySession
	for rows.Next() {
		var s StudySession
		if err := rows.Scan(&s.ID, &s.GroupID, &s.CreatedAt, &s.StudyActivityID); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// GetStudySession returns a study session by ID
func GetStudySession(id int64) (*StudySession, error) {
	var s StudySession
	err := DB.QueryRow(`
		SELECT id, group_id, created_at, study_activity_id 
		FROM study_sessions 
		WHERE id = ?
	`, id).Scan(&s.ID, &s.GroupID, &s.CreatedAt, &s.StudyActivityID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetStudySessionWords returns all words reviewed in a study session
func GetStudySessionWords(sessionID int64) ([]Word, error) {
	query := `
		SELECT DISTINCT w.id, w.japanese, w.romaji, w.english, w.parts
		FROM words w
		JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
	`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Japanese, &w.Romaji, &w.English, &w.Parts); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, nil
}

// CreateStudySession creates a new study session
func CreateStudySession(groupID int64) (*StudySession, error) {
	result, err := DB.Exec(`
		INSERT INTO study_sessions (group_id)
		VALUES (?)
	`, groupID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetStudySession(id)
}

// LinkStudyActivityToSession links a study activity to a session
func LinkStudyActivityToSession(activityID, sessionID int64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update study_activities
	_, err = tx.Exec(`
		UPDATE study_activities 
		SET study_session_id = ? 
		WHERE id = ?
	`, sessionID, activityID)
	if err != nil {
		return err
	}

	// Update study_sessions
	_, err = tx.Exec(`
		UPDATE study_sessions 
		SET study_activity_id = ? 
		WHERE id = ?
	`, activityID, sessionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetStudySessionStats returns statistics for a study session
func GetStudySessionStats(sessionID int64) (map[string]int, error) {
	var totalReviews, correctReviews int

	err := DB.QueryRow(`
		SELECT COUNT(*), SUM(CASE WHEN correct THEN 1 ELSE 0 END)
		FROM word_review_items
		WHERE study_session_id = ?
	`, sessionID).Scan(&totalReviews, &correctReviews)
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total_reviews":    totalReviews,
		"correct_reviews":  correctReviews,
		"accuracy_percent": (correctReviews * 100) / totalReviews,
	}, nil
}

// DeleteStudySession deletes a study session and its relationships
func DeleteStudySession(id int64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update study_activities to remove session link
	_, err = tx.Exec(`
		UPDATE study_activities 
		SET study_session_id = NULL 
		WHERE study_session_id = ?
	`, id)
	if err != nil {
		return err
	}

	// Delete word review items
	_, err = tx.Exec("DELETE FROM word_review_items WHERE study_session_id = ?", id)
	if err != nil {
		return err
	}

	// Delete the session
	_, err = tx.Exec("DELETE FROM study_sessions WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
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
		var a StudyActivity
		if err := rows.Scan(&a.ID, &a.StudySessionID, &a.GroupID, &a.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}