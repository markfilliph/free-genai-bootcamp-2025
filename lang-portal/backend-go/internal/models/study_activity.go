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