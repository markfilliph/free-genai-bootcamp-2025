package models

import (
	"database/sql"
	"time"
)

type LastStudySession struct {
	ID              int       `json:"id"`
	GroupID         int       `json:"group_id"`
	GroupName       string    `json:"group_name"`
	StudyActivityID int       `json:"study_activity_id"`
	ActivityName    string    `json:"activity_name"`
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CorrectCount    int       `json:"correct_count"`
	TotalCount      int       `json:"total_count"`
}

type StudyProgress struct {
	Date         string `json:"date"`
	CorrectCount int    `json:"correct_count"`
	TotalCount   int    `json:"total_count"`
}

type QuickStats struct {
	TotalWords      int     `json:"total_words"`
	TotalGroups     int     `json:"total_groups"`
	TotalSessions   int     `json:"total_sessions"`
	CorrectRate     float64 `json:"correct_rate"`
	StudiedWords    int     `json:"studied_words"`
	UnstudiedWords  int     `json:"unstudied_words"`
}

// GetLastStudySession retrieves the most recent study session with stats
func GetLastStudySession(db *sql.DB) (*LastStudySession, error) {
	var session LastStudySession
	err := db.QueryRow(`
		SELECT 
			ss.id,
			ss.group_id,
			g.name as group_name,
			ss.study_activity_id,
			sa.name as activity_name,
			ss.created_at,
			ss.completed_at,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(wri.word_id) as total_count
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY ss.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`).Scan(
		&session.ID,
		&session.GroupID,
		&session.GroupName,
		&session.StudyActivityID,
		&session.ActivityName,
		&session.CreatedAt,
		&session.CompletedAt,
		&session.CorrectCount,
		&session.TotalCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetStudyProgress retrieves study progress for the last 7 days
func GetStudyProgress(db *sql.DB) ([]StudyProgress, error) {
	rows, err := db.Query(`
		WITH RECURSIVE dates(date) AS (
			SELECT date('now', '-6 days')
			UNION ALL
			SELECT date(date, '+1 day')
			FROM dates
			WHERE date < date('now')
		)
		SELECT 
			dates.date,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(wri.word_id) as total_count
		FROM dates
		LEFT JOIN study_sessions ss ON date(ss.created_at) = dates.date
		LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
		GROUP BY dates.date
		ORDER BY dates.date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []StudyProgress
	for rows.Next() {
		var p StudyProgress
		if err := rows.Scan(&p.Date, &p.CorrectCount, &p.TotalCount); err != nil {
			return nil, err
		}
		progress = append(progress, p)
	}
	return progress, nil
}

// GetQuickStats retrieves quick statistics about words and study sessions
func GetQuickStats(db *sql.DB) (*QuickStats, error) {
	var stats QuickStats
	
	// Get total words and groups
	err := db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM words) as total_words,
			(SELECT COUNT(*) FROM groups) as total_groups,
			(SELECT COUNT(*) FROM study_sessions) as total_sessions
	`).Scan(&stats.TotalWords, &stats.TotalGroups, &stats.TotalSessions)
	if err != nil {
		return nil, err
	}

	// Get correct rate and studied/unstudied words
	err = db.QueryRow(`
		WITH word_stats AS (
			SELECT 
				COUNT(DISTINCT word_id) as studied_words,
				SUM(CASE WHEN correct THEN 1 ELSE 0 END) * 1.0 / COUNT(*) as correct_rate
			FROM word_review_items
		)
		SELECT 
			COALESCE(correct_rate, 0),
			COALESCE(studied_words, 0),
			(SELECT COUNT(*) FROM words) - COALESCE(studied_words, 0)
		FROM word_stats
	`).Scan(&stats.CorrectRate, &stats.StudiedWords, &stats.UnstudiedWords)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
