package models

import (
	"database/sql"
	"time"
)

type StudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	StudyActivityID int       `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// GetStudySession retrieves a study session by its ID
func GetStudySession(db *sql.DB, id int) (*StudySession, error) {
	var session StudySession
	query := `
		SELECT id, group_id, study_activity_id, created_at, completed_at 
		FROM study_sessions 
		WHERE id = $1`
	
	err := db.QueryRow(query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.CreatedAt,
		&session.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &session, nil
}
