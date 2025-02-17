package models

import (
	"database/sql"
	"time"
)

type StudyActivity struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

type StudySessionDetail struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int       `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
	ActivityName    string    `json:"activity_name"`
	ReviewItemCount int       `json:"review_items_count"`
}

type WordReviewItem struct {
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetStudyActivity retrieves a single study activity
func GetStudyActivity(db *sql.DB, id int) (*StudyActivity, error) {
	var sa StudyActivity
	err := db.QueryRow(`
		SELECT id, name, thumbnail_url, description 
		FROM study_activities 
		WHERE id = ?`,
		id).Scan(&sa.ID, &sa.Name, &sa.ThumbnailURL, &sa.Description)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

// GetStudySessions retrieves study sessions for an activity
func GetStudySessions(db *sql.DB, activityID int, page, perPage int) ([]StudySessionDetail, int, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM study_sessions 
		WHERE study_activity_id = ?`,
		activityID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated sessions with details
	rows, err := db.Query(`
		SELECT 
			ss.id, ss.group_id, ss.created_at, ss.study_activity_id,
			g.name as group_name,
			sa.name as activity_name,
			COUNT(wri.word_id) as review_items_count
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		WHERE ss.study_activity_id = ?
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?`,
		activityID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []StudySessionDetail
	for rows.Next() {
		var s StudySessionDetail
		if err := rows.Scan(
			&s.ID, &s.GroupID, &s.CreatedAt, &s.StudyActivityID,
			&s.GroupName, &s.ActivityName, &s.ReviewItemCount,
		); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

// CreateStudySession creates a new study session
func CreateStudySession(db *sql.DB, groupID, activityID int) (*StudySessionDetail, error) {
	result, err := db.Exec(`
		INSERT INTO study_sessions (group_id, study_activity_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)`,
		groupID, activityID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &StudySessionDetail{
		ID:              int(id),
		GroupID:         groupID,
		StudyActivityID: activityID,
		CreatedAt:       time.Now(),
	}, nil
}

// AddWordReview adds a word review to a study session
func AddWordReview(db *sql.DB, sessionID, wordID int, correct bool) error {
	_, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		wordID, sessionID, correct)
	return err
}
